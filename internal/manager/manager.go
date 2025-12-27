package manager

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/syscc/Emby-Go/internal/db"
	"github.com/syscc/Emby-Go/internal/util/logs"
	"gopkg.in/yaml.v3"
)

type serverProc struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

var (
	mu       sync.Mutex
	procs    = make(map[uint]*serverProc)
	DataRoot = "." // set by main
)

func writeConfig(_ string, s db.EmbyServer) (string, error) {
	var root map[string]any
	root = map[string]any{}
	gc, _ := db.GetGlobalConfig()

	emby := getMap(root, "emby")
	openlist := getMap(root, "openlist")
	strm := getMap(emby, "strm")
	path := getMap(root, "path")
	cache := getMap(root, "cache")
	vp := getMap(root, "video-preview")
	log := getMap(root, "log")
	ssl := getMap(root, "ssl")
	ltg := getMap(openlist, "local-tree-gen")

	// Emby Config
	emby["name"] = s.Name
	emby["host"] = s.EmbyHost
	if s.EmbyToken != "" {
		emby["token"] = s.EmbyToken
	}
	emby["episodes-unplay-prior"] = gc.EpisodesUnplayPrior
	emby["resort-random-items"] = gc.ResortRandomItems
	emby["proxy-error-strategy"] = gc.ProxyErrorStrategy
	emby["images-quality"] = gc.ImagesQuality
	emby["download-strategy"] = gc.DownloadStrategy
	if s.MountPath != "" {
		paths := splitMounts(s.MountPath)
		if len(paths) > 0 {
			emby["mount-path"] = strings.Join(paths, ",")
		}
	}
	if s.LocalMediaRoot != "" {
		emby["local-media-root"] = s.LocalMediaRoot
	}
	if s.DirectLinkCacheExpired != "" {
		emby["dl-cache-time"] = s.DirectLinkCacheExpired
	}

	// Strm Config
	strm["internal-redirect-enable"] = s.InternalRedirectEnable
	// StrmPathMap from global config
	if gc.StrmPathMap != "" {
		strm["path-map"] = strings.Split(gc.StrmPathMap, "\n")
	}

	// Openlist Config
	if s.OpenlistHost != "" {
		openlist["host"] = s.OpenlistHost
	}
	if s.OpenlistToken != "" {
		openlist["token"] = s.OpenlistToken
	}

	// Local Tree Gen Config (Missing in old manager)
	ltg["enable"] = gc.LTGEnable
	ltg["ffmpeg-enable"] = gc.LTGFFmpegEnable
	ltg["virtual-containers"] = gc.LTGVirtualContainers
	ltg["strm-containers"] = gc.LTGStrmContainers
	ltg["music-containers"] = gc.LTGMusicContainers
	ltg["auto-remove-max-count"] = gc.LTGAutoRemoveMaxCount
	ltg["refresh-interval"] = gc.LTGRefreshInterval
	if gc.LTGScanPrefixes != "" {
		ltg["scan-prefixes"] = strings.Split(gc.LTGScanPrefixes, "\n")
	}
	ltg["ignore-containers"] = gc.LTGIgnoreContainers
	ltg["threads"] = gc.LTGThreads

	// Cache Config
	cache["enable"] = gc.CacheEnable
	cache["expired"] = gc.CacheExpired
	if gc.CacheWhiteList != "" {
		cache["whitelist"] = strings.Split(gc.CacheWhiteList, "\n")
	}

	// Video Preview Config
	vp["enable"] = gc.VideoPreviewEnable
	if gc.VideoPreviewContainers != "" {
		vp["containers"] = strings.Split(gc.VideoPreviewContainers, ",")
	}
	if gc.VideoPreviewIgnoreTemplateIds != "" {
		vp["ignore-template-ids"] = strings.Split(gc.VideoPreviewIgnoreTemplateIds, ",")
	}

	// Path Config
	if gc.PathEmby2Openlist != "" {
		path["emby2openlist"] = strings.Split(gc.PathEmby2Openlist, "\n")
	}
	if s.MountPath != "" {
		for _, p := range splitMounts(s.MountPath) {
			// Ensure the mount path is mapped to root "/"
			ensurePathMap(path, "emby2openlist", fmt.Sprintf("%s:/", strings.TrimSpace(p)))
		}
	}

	// Log Config
	log["disable-color"] = gc.LogDisableColor

	// SSL Config (Missing in old manager)
	ssl["enable"] = gc.SslEnable
	ssl["single-port"] = gc.SslSinglePort
	ssl["key"] = gc.SslKey
	ssl["crt"] = gc.SslCrt

	data, err := yaml.Marshal(root)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(DataRoot, "servers", s.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	p := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return "", err
	}
	return p, nil
}

func getMap(m map[string]any, k string) map[string]any {
	v, ok := m[k]
	if !ok {
		v = map[string]any{}
		m[k] = v
	}
	mv, _ := v.(map[string]any)
	if mv == nil {
		mv = map[string]any{}
		m[k] = mv
	}
	return mv
}

func splitMounts(s string) []string {
	arr := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ';'
	})
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		e = strings.TrimSpace(e)
		if e != "" {
			out = append(out, e)
		}
	}
	return out
}

func ensurePathMap(m map[string]any, k string, entry string) {
	v, ok := m[k]
	if !ok || v == nil {
		m[k] = []any{entry}
		return
	}
	switch arr := v.(type) {
	case []any:
		for _, e := range arr {
			if s, ok := e.(string); ok && s == entry {
				return
			}
		}
		m[k] = append(arr, entry)
	case []string:
		for _, s := range arr {
			if s == entry {
				return
			}
		}
		ss := append(arr, entry)
		xx := make([]any, len(ss))
		for i, s := range ss {
			xx[i] = s
		}
		m[k] = xx
	default:
		m[k] = []any{entry}
	}
}

func LoadAll() error {
	list, err := db.GetServers()
	if err != nil {
		return err
	}
	for _, s := range list {
		go Start(s)
	}
	return nil
}

func Start(s db.EmbyServer) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := procs[s.ID]; ok {
		return
	}
	if s.DisableProxy {
		return
	}
	base := DataRoot
	cfgPath, err := writeConfig(base, s)
	if err != nil {
		logs.Error("写入配置失败: %v", err)
		return
	}

	exe, err := os.Executable()
	if err != nil {
		logs.Error("获取可执行文件路径失败: %v", err)
		return
	}

	// Use -kernel-only flag and point to the generated config
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, exe, "-dr", base, "-kernel-only", "-config", cfgPath, "-http-port", fmt.Sprintf("%d", s.HTTPPort), "-https-port", fmt.Sprintf("%d", s.HTTPPort-1))
	cmd.Env = os.Environ()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		logs.Error("启动内核失败: %v", err)
		cancel()
		return
	}

	go captureLogs(s, stdout)
	go captureLogs(s, stderr)

	procs[s.ID] = &serverProc{cmd: cmd, cancel: cancel}
	logs.Info("已启动 %s, 端口: %d, 配置: %s", s.Name, s.HTTPPort, cfgPath)
}

func captureLogs(s db.EmbyServer, r io.ReadCloser) {
	defer r.Close()
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if len(line) > 0 {
			level := "INFO"
			if strings.Contains(line, "[ERROR]") {
				level = "ERROR"
			} else if strings.Contains(line, "[WARN]") {
				level = "WARN"
			} else if strings.Contains(line, "[SUCCESS]") {
				level = "SUCCESS"
			}
			switch level {
			case "ERROR":
				logs.Error("[%s] %s", s.Name, strings.TrimSpace(line))
			case "WARN":
				logs.Warn("[%s] %s", s.Name, strings.TrimSpace(line))
			case "SUCCESS":
				logs.Success("[%s] %s", s.Name, strings.TrimSpace(line))
			default:
				logs.Info("[%s] %s", s.Name, strings.TrimSpace(line))
			}
		}
		if err != nil {
			return
		}
	}
}

func Stop(id uint) {
	mu.Lock()
	defer mu.Unlock()
	if sp, ok := procs[id]; ok {
		sp.cancel()
		done := make(chan struct{})
		go func() {
			sp.cmd.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
		}
		delete(procs, id)
		logs.Info("已停止服务: %d", id)
	}
}

func Restart(id uint) error {
	Stop(id)
	list, err := db.GetServers()
	if err != nil {
		return err
	}
	for _, s := range list {
		if s.ID == id {
			go Start(s)
			break
		}
	}
	return nil
}
