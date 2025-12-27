package cache

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/syscc/Emby-Go/internal/constant"
	"github.com/syscc/Emby-Go/internal/util/encrypts"
	"github.com/syscc/Emby-Go/internal/util/https"
	"github.com/syscc/Emby-Go/internal/util/logs"
	"github.com/syscc/Emby-Go/internal/util/logs/colors"
	"github.com/syscc/Emby-Go/internal/util/strs"
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
	"Sec-Fetch-Dest": {}, "Sec-Fetch-Mode": {}, "Sec-Fetch-Site": {},
	"Sec-Ch-Ua": {}, "Sec-Ch-Ua-Mobile": {}, "Sec-Ch-Ua-Platform": {},
	"Content-Type": {}, "Content-Length": {},

	// Emby Auth & Device (忽略这些以实现多用户/多设备共享缓存)
	"X-Emby-Token": {}, "X-Emby-Device-Id": {}, "X-Emby-Device-Name": {},
	"X-Emby-Client": {}, "X-Emby-Client-Version": {},

	// StreamMusic
	"X-Streammusic-Audioid": {}, "X-Streammusic-Savepath": {},

	// CDN & Proxy Headers (忽略 CDN 动态头)
	"X-Forwarded-For": {}, "X-Real-IP": {}, "X-Real-Ip": {}, "Forwarded": {}, "Client-IP": {},
	"True-Client-IP": {}, "CF-Connecting-IP": {}, "X-Cluster-Client-IP": {},
	"Fastly-Client-IP": {}, "X-Client-IP": {}, "X-ProxyUser-IP": {},
	"Via": {}, "Forwarded-For": {}, "X-From-Cdn": {},
	"Ali-Swift-Log-Host": {}, "Eagleid": {}, "X-Amz-Cf-Id": {}, "X-Request-Id": {},
	"X-Via": {}, "X-Ca-Request-Id": {}, "X-Nginx-Proxy": {},
}

// CacheableRouteMarker 缓存白名单
// 只有匹配上正则表达式的路由才会被缓存
func CacheableRouteMarker() gin.HandlerFunc {
	cacheablePatterns := []*regexp.Regexp{
		regexp.MustCompile(constant.Reg_PlaybackInfo),
		regexp.MustCompile(constant.Reg_VideoSubtitles),
		regexp.MustCompile(constant.Reg_ResourceStream),
		regexp.MustCompile(constant.Reg_ResourceOriginal), // 缓存 original 接口
		regexp.MustCompile(constant.Reg_ItemDownload),
		regexp.MustCompile(constant.Reg_ItemSyncDownload),
		regexp.MustCompile(constant.Reg_UserItemsRandomWithLimit),
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
		header := c.Writer.Header()
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
	
	c.Request.URL.RawQuery = q.Encode()
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
	header := strings.Builder{}
	for key, values := range c.Request.Header {
		if _, ok := CacheKeyIgnoreParams[key]; ok {
			continue
		}
		header.WriteString(key)
		header.WriteString("=")
		header.WriteString(strings.Join(values, "|"))
		header.WriteString(";")
	}

	headerStr := header.String()
	preEnc := strs.Sort(c.Request.URL.RawQuery + body + headerStr)
	if headerStr != "" {
		logs.Tip("headers to encode cacheKey: %s", colors.ToYellow(headerStr))
	}

	// 为防止字典排序后, 不同的 uri 冲突, 这里在排序完的字符串前再加上原始的 uri
	uriNoArgs := urls.ReplaceAll(
		uri,
		"?"+c.Request.URL.RawQuery, "",
		c.Request.URL.RawQuery, "",
	)

	hash := encrypts.Md5Hash(method + uriNoArgs + preEnc)
	return hash, nil
}
