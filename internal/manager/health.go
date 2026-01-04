package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/syscc/Emby-Go/internal/db"
	"github.com/syscc/Emby-Go/internal/util/logs"
)

func StartHealthMonitor() {
	go func() {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case <-t.C:
				checkAll()
			}
		}
	}()
}

func checkAll() {
	servers, err := db.GetServers()
	if err != nil {
		return
	}
	for _, s := range servers {
		if s.DisableProxy {
			continue
		}
		if unhealthy(s) {
			logs.Warn("健康检查: %s 返回 50x", s.Name)
			fail := true
			for i := 0; i < 2; i++ {
				time.Sleep(2 * time.Second)
				if !unhealthy(s) {
					fail = false
					break
				}
			}
			if fail {
				logs.Warn("健康检查: 重启服务 %s", s.Name)
				_ = Restart(s.ID)
				time.Sleep(5 * time.Second)
				if unhealthy(s) {
					logs.Error("健康检查: 服务 %s 重启后仍异常，发送通知", s.Name)
					_ = SendWebhookAll("Go-Emby 健康检查告警", "服务 "+s.Name+" 状态异常，重启后仍异常")
				} else {
					logs.Success("健康检查: 服务 %s 重启后恢复正常", s.Name)
				}
			}
		}
	}
}

func unhealthy(s db.EmbyServer) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + strconvI(s.HTTPPort) + "/")
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 500 && resp.StatusCode < 600
}

func SendWebhook(title, content string) error {
	g, err := db.GetGlobalConfig()
	if err != nil {
		return err
	}
	if !g.NotifyEnable || g.NotifyUrl == "" {
		return fmt.Errorf("notify disabled or url empty")
	}
	method := g.NotifyMethod
	if method == "" {
		method = "POST"
	}
	ct := g.NotifyContentType
	if ct == "" {
		ct = "application/json"
	}
	tk := g.NotifyTitleKey
	if tk == "" {
		tk = "title"
	}
	ck := g.NotifyContentKey
	if ck == "" {
		ck = "text"
	}
	var body []byte
	var req *http.Request
	switch ct {
	case "application/x-www-form-urlencoded":
		data := url.Values{}
		data.Set(tk, title)
		data.Set(ck, content)
		body = []byte(data.Encode())
		req, _ = http.NewRequest(method, g.NotifyUrl, bytes.NewBuffer(body))
	default:
		m := map[string]string{tk: title, ck: content}
		body, _ = json.Marshal(m)
		req, _ = http.NewRequest(method, g.NotifyUrl, bytes.NewBuffer(body))
	}
	req.Header.Set("Content-Type", ct)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logs.Warn("Webhook 发送失败: %v", err)
		return err
	}
	resp.Body.Close()
	return nil
}

func SendWebhookAll(title, content string) error {
	// Try list first
	list, err := db.GetNotifies()
	if err == nil && len(list) > 0 {
		for _, n := range list {
			if !n.Enable || n.Url == "" {
				continue
			}
			method := n.Method
			if method == "" {
				method = "POST"
			}
			ct := n.ContentType
			if ct == "" {
				ct = "application/json"
			}
			tk := n.TitleKey
			if tk == "" {
				tk = "title"
			}
			ck := n.ContentKey
			if ck == "" {
				ck = "text"
			}
			var body []byte
			var req *http.Request
			switch ct {
			case "application/x-www-form-urlencoded":
				data := url.Values{}
				data.Set(tk, title)
				data.Set(ck, content)
				body = []byte(data.Encode())
				req, _ = http.NewRequest(method, n.Url, bytes.NewBuffer(body))
			default:
				m := map[string]string{tk: title, ck: content}
				body, _ = json.Marshal(m)
				req, _ = http.NewRequest(method, n.Url, bytes.NewBuffer(body))
			}
			req.Header.Set("Content-Type", ct)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				logs.Warn("Webhook 发送失败: %v", err)
				continue
			}
			resp.Body.Close()
		}
		return nil
	}
	// Fallback single config
	return SendWebhook(title, content)
}

func strconvI(i int) string {
	return fmt.Sprintf("%d", i)
}
