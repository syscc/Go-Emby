package webui

import (
	"bufio"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syscc/Emby-Go/internal/db"
	"github.com/syscc/Emby-Go/internal/manager"
)

func trustedProxies() []string {
	val := strings.TrimSpace(os.Getenv("trusted_proxies"))
	if val != "" {
		parts := strings.Split(val, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return []string{
		"127.0.0.1/32", "::1/128",
		"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
	}
}

//go:embed static/*
var staticFS embed.FS

func Start(port int) {
	r := gin.Default()
	_ = r.SetTrustedProxies(trustedProxies())

	api := r.Group("/api")
	{
		api.GET("/check-init", func(c *gin.Context) {
			c.JSON(200, gin.H{"initialized": db.CheckInit()})
		})

		api.POST("/setup", func(c *gin.Context) {
			var form struct {
				Username string
				Password string
			}
			if err := c.ShouldBindJSON(&form); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.CreateUser(form.Username, form.Password); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})

		api.POST("/login", func(c *gin.Context) {
			var form struct {
				Username string
				Password string
			}
			if err := c.ShouldBindJSON(&form); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			user, err := db.GetUser(form.Username, form.Password)
			if err != nil {
				c.JSON(401, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(200, gin.H{"token": "dummy-token-" + strconv.Itoa(int(user.ID))})
		})

		// Protected routes
		auth := api.Group("/", func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.AbortWithStatus(401)
				return
			}
			c.Next()
		})

		auth.GET("/config", func(c *gin.Context) {
			wd, _ := os.Getwd()
			fp := filepath.Join(wd, "config.yml")
			b, err := os.ReadFile(fp)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"path": fp, "content": string(b)})
		})

		auth.PUT("/config", func(c *gin.Context) {
			var body struct {
				Content string
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			wd, _ := os.Getwd()
			fp := filepath.Join(wd, "config.yml")
			if err := os.WriteFile(fp, []byte(body.Content), 0o644); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})

		auth.GET("/global-config", func(c *gin.Context) {
			g, err := db.GetGlobalConfig()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, g)
		})
		auth.PUT("/global-config", func(c *gin.Context) {
			var g db.GlobalConfig
			if err := c.ShouldBindJSON(&g); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.UpdateGlobalConfig(&g); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			servers, _ := db.GetServers()
			for _, s := range servers {
				_ = manager.Restart(s.ID)
			}
			c.Status(200)
		})
		auth.GET("/notification", func(c *gin.Context) {
			g, err := db.GetGlobalConfig()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{
				"Enable":      g.NotifyEnable,
				"Url":         g.NotifyUrl,
				"Method":      g.NotifyMethod,
				"ContentType": g.NotifyContentType,
				"TitleKey":    g.NotifyTitleKey,
				"ContentKey":  g.NotifyContentKey,
			})
		})
		auth.PUT("/notification", func(c *gin.Context) {
			var body struct {
				Enable      bool
				Url         string
				Method      string
				ContentType string
				TitleKey    string
				ContentKey  string
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			g, err := db.GetGlobalConfig()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			g.NotifyEnable = body.Enable
			g.NotifyUrl = body.Url
			g.NotifyMethod = body.Method
			g.NotifyContentType = body.ContentType
			g.NotifyTitleKey = body.TitleKey
			g.NotifyContentKey = body.ContentKey
			if err := db.UpdateGlobalConfig(&g); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
		auth.POST("/notification/test", func(c *gin.Context) {
			var body struct {
				Title string
				Text  string
				// 兼容旧字段
				Content string
			}
			_ = c.ShouldBindJSON(&body)
			if body.Title == "" {
				body.Title = "Go-Emby Test"
			}
			msg := body.Text
			if msg == "" {
				msg = body.Content
			}
			if msg == "" {
				msg = "这是一条由您自己发送的Go-Emby测试消息，当你看到这条消息，说明你的配置是正确可用的。"
			}
			g, err := db.GetGlobalConfig()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if !g.NotifyEnable || strings.TrimSpace(g.NotifyUrl) == "" {
				c.JSON(400, gin.H{"error": "通知未启用或请求地址为空"})
				return
			}
			if err := manager.SendWebhook(body.Title, msg); err != nil {
				c.JSON(502, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
		auth.GET("/servers", func(c *gin.Context) {
			servers, _ := db.GetServers()
			c.JSON(200, servers)
		})

		// Notifications list CRUD
		auth.GET("/notifications", func(c *gin.Context) {
			list, err := db.GetNotifies()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, list)
		})
		auth.POST("/notifications", func(c *gin.Context) {
			var n db.Notify
			if err := c.ShouldBindJSON(&n); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.AddNotify(&n); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
		auth.PUT("/notifications/:id", func(c *gin.Context) {
			var n db.Notify
			if err := c.ShouldBindJSON(&n); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			id, _ := strconv.Atoi(c.Param("id"))
			n.ID = uint(id)
			if err := db.UpdateNotify(&n); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
		auth.DELETE("/notifications/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			if err := db.DeleteNotify(uint(id)); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
		auth.POST("/notifications/test", func(c *gin.Context) {
			var body db.Notify
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if !body.Enable || strings.TrimSpace(body.Url) == "" {
				c.JSON(400, gin.H{"error": "通知未启用或请求地址为空"})
				return
			}
			title := "Go-Emby Test"
			content := "这是一条由您自己发送的Go-Emby测试消息，当你看到这条消息，说明你的配置是正确可用的。"
			// Build request with provided body
			method := body.Method
			if method == "" {
				method = "POST"
			}
			ct := body.ContentType
			if ct == "" {
				ct = "application/json"
			}
			tk := body.TitleKey
			if tk == "" {
				tk = "title"
			}
			ck := body.ContentKey
			if ck == "" {
				ck = "text"
			}
			var req *http.Request
			switch ct {
			case "application/x-www-form-urlencoded":
				data := url.Values{}
				data.Set(tk, title)
				data.Set(ck, content)
				req, _ = http.NewRequest(method, body.Url, strings.NewReader(data.Encode()))
			default:
				m := map[string]string{tk: title, ck: content}
				b, _ := json.Marshal(m)
				req, _ = http.NewRequest(method, body.Url, strings.NewReader(string(b)))
			}
			req.Header.Set("Content-Type", ct)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				if resp != nil {
					resp.Body.Close()
				}
				c.JSON(502, gin.H{"error": "发送失败"})
				return
			}
			resp.Body.Close()
			c.Status(200)
		})

		auth.POST("/servers", func(c *gin.Context) {
			var s db.EmbyServer
			if err := c.ShouldBindJSON(&s); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.AddServer(&s); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			manager.Start(s)
			c.Status(200)
		})

		auth.PUT("/servers/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var s db.EmbyServer
			if err := c.ShouldBindJSON(&s); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			s.ID = uint(id)
			if err := db.UpdateServer(&s); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			manager.Restart(uint(id))
			c.Status(200)
		})

		auth.DELETE("/servers/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			manager.Stop(uint(id))
			db.DeleteServer(uint(id))
			c.Status(200)
		})

		auth.GET("/logs", func(c *gin.Context) {
			type LogLine struct {
				Level   string `json:"Level"`
				Message string `json:"Message"`
			}
			var list []LogLine
			fp := filepath.Join(manager.DataRoot, "log", "go-emby_"+time.Now().Format("2006-01-02")+".log")
			f, err := os.Open(fp)
			if err == nil {
				defer f.Close()
				var lines []string
				sc := bufio.NewScanner(f)
				for sc.Scan() {
					lines = append(lines, sc.Text())
				}
				// reverse latest first
				for i := len(lines) - 1; i >= 0 && len(list) < 500; i-- {
					parts := strings.SplitN(lines[i], "\t", 3)
					if len(parts) < 3 {
						continue
					}
					list = append(list, LogLine{
						Level:   parts[1],
						Message: parts[2],
						// TODO: add origin
					})
				}
			}
			c.JSON(200, list)
		})

		auth.POST("/user/password", func(c *gin.Context) {
			var form struct {
				Username    string
				NewPassword string
			}
			if err := c.ShouldBindJSON(&form); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.UpdatePassword(form.Username, form.NewPassword); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.Status(200)
		})
	}

	// Static files
	sub, _ := fs.Sub(staticFS, "static")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(sub))))
	// r.NoRoute(gin.WrapH(http.FileServer(http.Dir("internal/webui/static"))))

	r.Run(":" + strconv.Itoa(port))
}
