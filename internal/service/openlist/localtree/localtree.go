package localtree

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/syscc/Emby-Go/internal/config"
	"github.com/syscc/Emby-Go/internal/util/logs/colors"
)

// DirName 存放目录树的本地目录名称
const DirName = "openlist-local-tree"

// Init 根据配置文件, 初始化本地目录树
func Init() error {
	// 判断配置是否开启
	if !config.C.Openlist.LocalTreeGen.Enable {
		return nil
	}

	dirAbs := filepath.Join(config.BasePath, DirName)

	s := NewSynchronizer(dirAbs, 30)
	go startSync(s)

	return nil
}

// startSync 立即同步一次目录树, 并开始定时扫描同步变更
func startSync(s *Synchronizer) {
	doSync := func() {
		logf(colors.Blue, "开始同步")
		start := time.Now()
		total, added, deleted, err := s.Sync()
		if err != nil {
			logf(colors.Red, "同步失败: %v", err)
			return
		}
		logf(colors.Green, "同步完成, 总数: %d, 新增: %d, 删除: %d, 耗时: %v", total, added, deleted, time.Since(start))
	}
	doSync()

	d := time.Minute * time.Duration(config.C.Openlist.LocalTreeGen.RefreshInterval)
	timer := time.NewTicker(d)
	for range timer.C {
		doSync()
	}
}

// logf 带上前缀的日志输出
func logf(c colors.C, format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log.Println(colors.WrapColor(c, "[openlist 目录树]: "+s))
}
