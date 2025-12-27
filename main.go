package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/syscc/Emby-Go/internal/config"
	"github.com/syscc/Emby-Go/internal/constant"
	"github.com/syscc/Emby-Go/internal/db"
	"github.com/syscc/Emby-Go/internal/manager"
	"github.com/syscc/Emby-Go/internal/service/openlist/localtree"
	"github.com/syscc/Emby-Go/internal/util/logs"
	"github.com/syscc/Emby-Go/internal/util/logs/colors"
	"github.com/syscc/Emby-Go/internal/web"
	"github.com/syscc/Emby-Go/internal/web/webport"
	"github.com/syscc/Emby-Go/internal/webui"
	"github.com/gin-gonic/gin"
)

var ginMode = gin.DebugMode

func main() {
	go func() { http.ListenAndServe(":60360", nil) }()

	// Flags
	printVersion := flag.Bool("version", false, "查看程序版本")
	dr := flag.String("dr", "/app", "程序数据根目录") // Default to /app for docker environment
	wp := flag.Int("p", 8090, "WebUI 管理后台端口")

	// Kernel flags
	kernelOnly := flag.Bool("kernel-only", false, "仅启动内核")
	configPath := flag.String("config", "", "配置文件路径 (仅内核模式)")
	httpPort := flag.Int("http-port", 8095, "HTTP 端口 (仅内核模式)")
	httpsPort := flag.Int("https-port", 8094, "HTTPS 端口 (仅内核模式)")

	flag.Parse()

	// Load env files
	loadEnv()

	// Apply env vars if flags are default
	if *wp == 8090 {
		if v := os.Getenv("webui"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				*wp = p
			}
		}
	}
	if *dr == "./app" {
		if v := os.Getenv("dr"); v != "" {
			*dr = v
		}
	}
	// kernel-only is bool, explicit flag usually overrides env, but if flag is false (default)
	// and env is true, we should set it?
	if !*kernelOnly {
		if v := os.Getenv("kernel_only"); v == "true" || v == "1" {
			*kernelOnly = true
		}
	}
	if *configPath == "" {
		if v := os.Getenv("config"); v != "" {
			*configPath = v
		}
	}
	if *httpPort == 8095 {
		if v := os.Getenv("http_port"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				*httpPort = p
			}
		}
	}
	if *httpsPort == 8094 {
		if v := os.Getenv("https_port"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				*httpsPort = p
			}
		}
	}

	if *printVersion {
		fmt.Println(constant.CurrentVersion)
		os.Exit(0)
	}

	if *kernelOnly {
		// --- Kernel Mode ---
		if *configPath == "" {
			log.Fatal("kernel-only mode requires -config")
		}

		// Setup logs for kernel
		serverDir := filepath.Dir(*configPath)
		logDir := filepath.Join(serverDir, "log")
		_ = os.MkdirAll(logDir, 0o755)
		setupLogHook(logDir, "go-emby")

		logs.Info("Starting Kernel: %s", *configPath)

		// Set ports
		webport.HTTP = strconv.Itoa(*httpPort)
		webport.HTTPS = strconv.Itoa(*httpsPort)

		// Load config
		if err := config.ReadFromFile(*configPath); err != nil {
			log.Fatal(err)
		}

		printBanner()

		logs.Info("正在初始化本地目录树模块...")
		if err := localtree.Init(); err != nil {
			log.Fatal(colors.ToRed(err.Error()))
		}

		logs.Info("正在启动服务...")
		gin.SetMode(ginMode)
		if err := web.Listen(); err != nil {
			log.Fatal(colors.ToRed(err.Error()))
		}
	} else {
		// --- Manager Mode ---
		manager.DataRoot = *dr

		printBanner()

		// Setup global log
		logDir := filepath.Join(*dr, "log")
		_ = os.MkdirAll(logDir, 0o755)
		setupLogHook(logDir, "go-emby")

		// Create directories
		_ = os.MkdirAll(filepath.Join(*dr, "custom-js"), 0o755)
		_ = os.MkdirAll(filepath.Join(*dr, "custom-css"), 0o755)
		_ = os.MkdirAll(filepath.Join(*dr, "lib"), 0o755)

		// Init DB
		dbPath := filepath.Join(*dr, "Emby-Go.db")
		logs.Info("正在初始化数据库: %s", dbPath)
		if err := db.Init(dbPath); err != nil {
			log.Fatalf("Init DB failed: %v", err)
		}

		// Start Manager (Load Proxies)
		logs.Info("正在加载代理服务...")
		if err := manager.LoadAll(); err != nil {
			log.Printf("Load proxies failed: %v", err)
		}

		// Start WebUI
		logs.Info("正在启动 WebUI 管理后台...")
		webui.Start(*wp)

		// Block forever
		select {}
	}
}

func setupLogHook(logDir, prefix string) {
	logs.OutputHook = func(level, msg string) {
		day := time.Now().Format("2006-01-02")
		fp := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", prefix, day))
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		defer f.Close()
		line := fmt.Sprintf("%s\t%s\t%s\n", time.Now().Format("2006/01/02 15:04:05"), level, msg)
		_, _ = f.WriteString(line)
	}
}

func printBanner() {
	fmt.Printf(colors.ToYellow(`
                                 _           ___                        _ _     _   
                                | |         |__ \                      | (_)   | |  
  __ _  ___ ______ ___ _ __ ___ | |__  _   _   ) |___  _ __   ___ _ __ | |_ ___| |_ 
 / _| |/ _ \______/ _ \ '_ | _ \| '_ \| | | | / // _ \| '_ \ / _ \ '_ \| | / __| __|
| (_| | (_) |    |  __/ | | | | | |_) | |_| |/ /| (_) | |_) |  __/ | | | | \__ \ |_ 
 \__, |\___/      \___|_| |_| |_|_.__/ \__, |____\___/| .__/ \___|_| |_|_|_|___/\__|
  __/ |                                 __/ |         | |                           
 |___/                                 |___/          |_|                           

 Repository: %s
    Version: %s
`), constant.RepoAddr, constant.CurrentVersion)
}

func loadEnv() {
	files := []string{".env", ".github/.env", "/app/.env"}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Don't overwrite existing env vars
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}
