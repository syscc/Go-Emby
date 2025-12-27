package logs

import (
	"fmt"
	"time"

	"github.com/syscc/Emby-Go/internal/util/logs/colors"
)

var OutputHook func(level, msg string)

func log(level, msg string) {
	if OutputHook != nil {
		OutputHook(level, msg)
	}
}

// Info 输出蓝色 Info 日志
func Info(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log("INFO", s)
	fmt.Println(now() + colors.ToBlue("[INFO] "+s))
}

// Success 输出绿色 Success 日志
func Success(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log("SUCCESS", s)
	fmt.Println(now() + colors.ToGreen("[SUCCESS] "+s))
}

// Warn 输出黄色 Warn 日志
func Warn(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log("WARN", s)
	fmt.Println(now() + colors.ToYellow("[WARN] "+s))
}

// Error 输出红色 Error 日志
func Error(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log("ERROR", s)
	fmt.Println(now() + colors.ToRed("[ERROR] "+s))
}

// Tip 输出灰色 Tip 日志
func Tip(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	// Tip doesn't have a level prefix in original, but for hook we can use TIP
	log("TIP", s)
	fmt.Println(now() + colors.ToGray(s))
}

// Progress 输出紫色 Progress 日志
func Progress(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	log("PROGRESS", s)
	fmt.Println(now() + colors.ToPurple(s))
}

// now 返回当前时间戳
func now() string {
	return time.Now().Format("2006-01-02 15:04:05") + " "
}
