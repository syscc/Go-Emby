package cache

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/syscc/Emby-Go/internal/constant"
	"github.com/syscc/Emby-Go/internal/util/encrypts"
	"github.com/syscc/Emby-Go/internal/util/https"
	"github.com/syscc/Emby-Go/internal/util/logs"
	"github.com/syscc/Emby-Go/internal/util/urls"

	"github.com/gin-gonic/gin"
)

// CacheKeyIgnoreParams 忽略的请求头或者参数
//
// 如果请求地址包含列表中的请求头或者参数, 则不参与 cacheKey 运算
var CacheKeyIgnoreParams = map[string]struct{}{
	// Fileball
	"StartTimeTicks": {}, "X-Playback-Session-Id": {},

	// Emby
	"PlaySessionId": {},

	// Common
	"Range": {}, "Host": {}, "Referrer": {}, "Connection": {},
	"Accept": {}, "Accept-Encoding": {}, "Accept-Language": {}, "Cache-Control": {},
	"Upgrade-Insecure-Requests": {}, "Referer": {}, "Origin": {}, "User-Agent": {},
	"Pragma": {}, "Priority": {}, "Cookie": {},

	// 忽略浏览器动态头
	"Sec-Fetch-Dest": {}, "Sec-Fetch-Mode": {}, "Sec-Fetch-Site": {},
	"Sec-Ch-Ua": {}, "Sec-Ch-Ua-Mobile": {}, "Sec-Ch-Ua-Platform": {},
	"Content-Type": {}, "Content-Length": {},

	// 忽略 Auth & Identity (实现多用户共享缓存)
	"X-Emby-Token": {}, "X-Emby-Authorization": {}, "Authorization": {},
	"X-Emby-Device-Id": {}, "X-Emby-Device-Name": {},

	// 忽略 IP & Proxy & CDN (实现任意 IP 共享缓存)
	"X-Emby-Client-IP": {}, "X-Real-IP": {}, "X-Forwarded-For": {},
	"X-Forwarded-Proto": {}, "X-Forwarded-Host": {}, "X-Forwarded-Port": {},
	"Via": {}, "X-Via": {},
	"Cf-Ray": {}, "Cf-Visitor": {}, "Cf-Connecting-Ip": {}, "Cf-Ipcountry": {},
	"Ali-Cdn-Real-Ip": {}, "Ali-Swift-Log-Host": {},
}

// CacheableRouteMarker 缓存白名单
// 只有匹配上正则表达式的路由才会被缓存
func CacheableRouteMarker() gin.HandlerFunc {
	cacheablePatterns := []*regexp.Regexp{
		// regexp.MustCompile(constant.Reg_PlaybackInfo),
		// regexp.MustCompile(constant.Reg_VideoSubtitles),
		regexp.MustCompile(constant.Reg_ResourceStream),
		regexp.MustCompile(constant.Reg_ResourceOriginal), // 缓存 original 接口
		// regexp.MustCompile(constant.Reg_ItemDownload),
		// regexp.MustCompile(constant.Reg_ItemSyncDownload),
		// regexp.MustCompile(constant.Reg_UserItemsRandomWithLimit),
	}

	return func(c *gin.Context) {
		for _, pattern := range cacheablePatterns {
			if pattern.MatchString(c.Request.RequestURI) {
				return
			}
		}
		c.Header(HeaderKeyExpired, "-1")
	}
}

// RequestCacher 请求缓存中间件
func RequestCacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1 判断请求是否需要缓存
		if c.Writer.Header().Get(HeaderKeyExpired) == "-1" {
			return
		}

		// 2 计算 cache key
		cacheKey, err := calcCacheKey(c)
		if err != nil {
			logs.Warn("cache key 计算异常: %v, 跳过缓存", err)
			// 如果没有调用 Abort, Gin 会自动继续调用处理器链
			return
		}

		// 3 尝试获取缓存
		if rc, ok := getCache(cacheKey); ok {
			// 调试日志：命中缓存
			logs.Tip("GetCache Hit: %s", cacheKey)
			if https.IsRedirectCode(rc.code) {
				// 适配重定向请求
				location := rc.header.header.Get("Location")
				c.Redirect(rc.code, location)
				logs.Success("直链缓存命中: %s", location)
			} else {
				c.Status(rc.code)
				https.CloneHeader(c.Writer, rc.header.header)
				c.Writer.Write(rc.body)
			}
			c.Abort()
			return
		}

		// 4 使用自定义的响应器
		customWriter := &respCacheWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = customWriter

		// 5 执行请求处理器
		c.Next()

		// 6 不缓存错误请求
		if https.IsErrorStatus(c.Writer.Status()) {
			return
		}

		// 7 刷新缓存
		// 如果缓存被中止 (通常是因为响应体过大), 则跳过
		if customWriter.aborted {
			logs.Warn("响应体过大, 跳过缓存: %s", cacheKey)
			return
		}

		header := c.Writer.Header()
		// logs.Tip("RespHeader Expired: '%s'", header.Get(HeaderKeyExpired)) // 调试日志
		respHeader := respHeader{
			expired:  header.Get(HeaderKeyExpired),
			space:    header.Get(HeaderKeySpace),
			spaceKey: header.Get(HeaderKeySpaceKey),
			header:   header.Clone(),
		}
		defer header.Del(HeaderKeyExpired)
		defer header.Del(HeaderKeySpace)
		defer header.Del(HeaderKeySpaceKey)

		go putCache(cacheKey, c, append([]byte(nil), customWriter.body.Bytes()...), respHeader)
	}
}

// Duration 将一个标准的时间转换成适用于缓存时间的字符串
func Duration(d time.Duration) string {
	expired := d.Milliseconds() + time.Now().UnixMilli()
	return fmt.Sprintf("%v", expired)
}

// WaitingForHandleChan 等待预缓存通道被处理完毕
func WaitingForHandleChan() {
	cacheHandleWaitGroup.Wait()
}

// calcCacheKey 计算缓存 key
//
// 计算方式: 取出 请求方法, 请求路径, 请求体, 请求头 转换成字符串之后字典排序,
// 再进行 Md5Hash
func calcCacheKey(c *gin.Context) (string, error) {
	method := c.Request.Method

	q := c.Request.URL.Query()
	for key := range CacheKeyIgnoreParams {
		q.Del(key)
	}
	// 额外忽略 Query 中的动态参数
	q.Del("api_key")
	q.Del("X-Emby-Token")
	q.Del("DeviceId")
	q.Del("reqformat")
	q.Del("static")
	q.Del("MediaSourceId")
	q.Del("PlaySessionId")
	q.Del("LiveStreamId")
	q.Del("StartTimeTicks")

	q.Del("AutoOpenLiveStream")
	q.Del("IsPlayback")

	// 忽略直链签名等动态参数
	q.Del("tempauth")
	q.Del("Translate")
	q.Del("UniqueId")
	q.Del("ApiVersion")

	// 忽略客户端 IP 相关参数 (关键)
	q.Del("X-Emby-Client-IP")
	q.Del("X-Real-IP")
	q.Del("X-Forwarded-For")

	// 再次尝试覆盖 RawQuery (关键)
	rawQuery := q.Encode()
	c.Request.URL.RawQuery = rawQuery
	uri := c.Request.URL.String()

	body := ""
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return "", fmt.Errorf("读取请求体失败: %v", err)
		}
		body = string(bodyBytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	// 1. 提取所有需要参与计算的 Header Keys
	// 为解决 CDN/不同客户端 导致的 Header 不一致问题，
	// 且我们的缓存主要针对 302 重定向链接 (不依赖 Header)，
	// 故在此处 强制忽略所有 Header，只根据 URL (Path + Filtered Query) 生成 Key。
	// 这样可以最大程度提高缓存命中率。
	headerStr := ""

	// 4. 构建用于 Hash 的完整字符串
	// c.Request.URL.RawQuery 已经是排序过的 (q.Encode() 会排序)
	// 为防止字典排序后, 不同的 uri 冲突, 这里在排序完的字符串前再加上原始的 uri
	uriNoArgs := urls.ReplaceAll(
		uri,
		"?"+c.Request.URL.RawQuery, "",
		c.Request.URL.RawQuery, "",
	)

	fullStringToHash := method + uriNoArgs + "?" + c.Request.URL.RawQuery + body + headerStr

	hash := encrypts.Md5Hash(fullStringToHash)
	return hash, nil
}
