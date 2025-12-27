package path

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/syscc/Emby-Go/internal/config"
	"github.com/syscc/Emby-Go/internal/service/openlist"
	"github.com/syscc/Emby-Go/internal/util/logs"
	"github.com/syscc/Emby-Go/internal/util/urls"
)

// OpenlistPathRes 路径转换结果
type OpenlistPathRes struct {

	// Success 转换是否成功
	Success bool

	// Path 转换后的路径
	Path string

	// Range 遍历所有 Openlist 根路径生成的子路径
	Range func() ([]string, error)
}

// Emby2Openlist Emby 资源路径转 Openlist 资源路径
func Emby2Openlist(embyPath string) OpenlistPathRes {
	pathRoutes := strings.Builder{}
	pathRoutes.WriteString("[")
	pathRoutes.WriteString("\n【原始路径】 => " + embyPath)

	embyPath = urls.Unescape(embyPath)
	pathRoutes.WriteString("\n\n【URL 解码】 => " + embyPath)

	embyPath = urls.TransferSlash(embyPath)
	pathRoutes.WriteString("\n\n【Windows 反斜杠转换】 => " + embyPath)

	embyMount := config.C.Emby.MountPath
	
	// 支持多个挂载路径，使用逗号或分号分隔
	mountPaths := strings.FieldsFunc(embyMount, func(r rune) bool {
		return r == ',' || r == ';'
	})

	var openlistFilePath string
	matched := false

	// 尝试匹配挂载路径
	for _, mount := range mountPaths {
		mount = strings.TrimSpace(mount)
		if mount == "" {
			continue
		}
		// 如果匹配成功，移除前缀
		if strings.HasPrefix(embyPath, mount) {
			openlistFilePath = strings.TrimPrefix(embyPath, mount)
			pathRoutes.WriteString("\n\n【移除 mount-path (" + mount + ")】 => " + openlistFilePath)
			matched = true
			break
		}
	}

	// 如果没有匹配到任何挂载路径，默认使用原逻辑（或视作未匹配）
	// 这里保持原有行为：如果未匹配，openlistFilePath 为空或者原值？
	// 原逻辑是 openlistFilePath := strings.TrimPrefix(embyPath, embyMount)
	// 如果 embyMount 只是单个路径且不匹配，openlistFilePath == embyPath
	
	if !matched {
		// 如果没有匹配到，此时 openlistFilePath 为空字符串
		// 为了保持兼容性，如果配置为空或者都没匹配上，我们假设它可能不需要移除或者就是原始路径
		// 但通常必须移除挂载点才能对应到 OpenList 的路径
		// 让我们暂且设为 embyPath，后续 mapEmby2Openlist 可能会处理
		openlistFilePath = embyPath
		pathRoutes.WriteString("\n\n【未匹配 mount-path】 => " + openlistFilePath)
	}

	if mapPath, ok := config.C.Path.MapEmby2Openlist(openlistFilePath); ok {
		openlistFilePath = mapPath
		pathRoutes.WriteString("\n\n【命中 emby2openlist 映射】 => " + openlistFilePath)
	}
	pathRoutes.WriteString("\n]")
	logs.Tip("embyPath 转换路径: %s", pathRoutes.String())

	rangeFunc := func() ([]string, error) {
		filePath, err := SplitFromSecondSlash(openlistFilePath)
		if err != nil {
			return nil, fmt.Errorf("openlistFilePath 解析异常: %s, error: %v", openlistFilePath, err)
		}

		res := openlist.FetchFsList("/", nil)
		if res.Code != http.StatusOK {
			return nil, fmt.Errorf("请求 openlist fs list 接口异常: %s", res.Msg)
		}

		paths := make([]string, 0, len(res.Data.Content))
		for _, c := range res.Data.Content {
			if !c.IsDir {
				continue
			}
			newPath := fmt.Sprintf("/%s%s", c.Name, filePath)
			paths = append(paths, newPath)
		}
		return paths, nil
	}

	return OpenlistPathRes{
		Success: true,
		Path:    openlistFilePath,
		Range:   rangeFunc,
	}
}

// SplitFromSecondSlash 找到给定字符串 str 中第二个 '/' 字符的位置
// 并以该位置为首字符切割剩余的子串返回
func SplitFromSecondSlash(str string) (string, error) {
	str = urls.TransferSlash(str)
	firstIdx := strings.Index(str, "/")
	if firstIdx == -1 {
		return "", fmt.Errorf("字符串不包含 /: %s", str)
	}

	secondIdx := strings.Index(str[firstIdx+1:], "/")
	if secondIdx == -1 {
		return "", fmt.Errorf("字符串只有单个 /: %s", str)
	}

	return str[secondIdx+firstIdx+1:], nil
}
