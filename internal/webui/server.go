package webui

import (
	"bufio"
	"embed"
	"io/fs"
	"net/http"
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
		auth.GET("/servers", func(c *gin.Context) {
			servers, _ := db.GetServers()
			c.JSON(200, servers)
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

	go r.Run(":" + strconv.Itoa(port))
}
