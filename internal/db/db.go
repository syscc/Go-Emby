package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string
}

type Notify struct {
 	ID           uint      `gorm:"primaryKey" json:"ID"`
	Name         string    `json:"Name"`
 	Enable       bool      `json:"Enable"`
 	Url          string    `json:"Url"`
 	Method       string    `json:"Method"`
 	ContentType  string    `json:"ContentType"`
 	TitleKey     string    `json:"TitleKey"`
 	ContentKey   string    `json:"ContentKey"`
 	CreatedAt    time.Time `json:"CreatedAt"`
 	UpdatedAt    time.Time `json:"UpdatedAt"`
 }

type EmbyServer struct {
	ID                     uint   `gorm:"primaryKey" json:"ID"`
	Name                   string `gorm:"uniqueIndex" json:"Name"`
	HTTPPort               int    `gorm:"uniqueIndex" json:"HTTPPort"`
	EmbyHost               string `json:"EmbyHost"`
	EmbyToken              string `json:"EmbyToken"`
	MountPath              string `json:"MountPath"`
	LocalMediaRoot         string `json:"LocalMediaRoot"` // Emby 本地媒体根目录
	OpenlistHost           string `json:"OpenlistHost"`
	OpenlistToken          string `json:"OpenlistToken"`
	InternalRedirectEnable bool   `json:"InternalRedirectEnable"`
	DirectLinkCacheExpired string `json:"DirectLinkCacheExpired"`
	DirectLinkCacheIgnore  string `json:"DirectLinkCacheIgnore"`
	DisableProxy           bool   `json:"DisableProxy"` // If true, only serve as config holder, don't start proxy
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type GlobalConfig struct {
	ID                            uint `gorm:"primaryKey"`
	EpisodesUnplayPrior           bool
	ResortRandomItems             bool
	ProxyErrorStrategy            string
	ImagesQuality                 int
	DownloadStrategy              string
	CacheEnable                   bool
	CacheExpired                  string
	VideoPreviewEnable            bool
	VideoPreviewContainers        string
	VideoPreviewIgnoreTemplateIds string
	PathEmby2Openlist             string
	LogDisableColor               bool
	StrmPathMap                   string
	CacheWhiteList                string
	LTGEnable                     bool
	LTGFFmpegEnable               bool
	LTGVirtualContainers          string
	LTGStrmContainers             string
	LTGMusicContainers            string
	LTGAutoRemoveMaxCount         int
	LTGRefreshInterval            int
	LTGScanPrefixes               string
	LTGIgnoreContainers           string
	LTGThreads                    int
	SslEnable                     bool
	SslSinglePort                 bool
	SslKey                        string
	SslCrt                        string
	NotifyEnable                  bool
	NotifyUrl                     string
	NotifyMethod                  string
	NotifyContentType             string
	NotifyTitleKey                string
	NotifyContentKey              string
}

func Init(path string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(&User{}, &EmbyServer{}, &GlobalConfig{}, &Notify{}); err != nil {
		return err
	}
	return ensureGlobalDefaults()
}

func CreateUser(username, password string) error {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count > 0 {
		return errors.New("admin user already exists")
	}
	return DB.Create(&User{Username: username, Password: hashMD5(password)}).Error
}

func GetUser(username, password string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	h := hashMD5(password)
	if user.Password == h {
		return &user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func CheckInit() bool {
	if DB == nil {
		return false
	}
	var count int64
	DB.Model(&User{}).Count(&count)
	return count > 0
}

func GetServers() ([]EmbyServer, error) {
	var servers []EmbyServer
	if err := DB.Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func AddServer(s *EmbyServer) error {
	return DB.Create(s).Error
}

func DeleteServer(id uint) error {
	return DB.Delete(&EmbyServer{}, id).Error
}

func UpdateServer(s *EmbyServer) error {
	return DB.Save(s).Error
}

func UpdatePassword(username, newPassword string) error {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	user.Password = hashMD5(newPassword)
	return DB.Save(&user).Error
}

func GetNotifies() ([]Notify, error) {
	var list []Notify
	if err := DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func AddNotify(n *Notify) error {
	return DB.Create(n).Error
}

func UpdateNotify(n *Notify) error {
	return DB.Save(n).Error
}

func DeleteNotify(id uint) error {
	return DB.Delete(&Notify{}, id).Error
}

func hashMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func ensureGlobalDefaults() error {
	var count int64
	DB.Model(&GlobalConfig{}).Count(&count)
	if count > 0 {
		return nil
	}
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "config.yml")
	b, err := os.ReadFile(fp)
	if err != nil {
		return DB.Create(&GlobalConfig{
			EpisodesUnplayPrior:           true,
			ResortRandomItems:             true,
			ProxyErrorStrategy:            "origin",
			ImagesQuality:                 100,
			DownloadStrategy:              "403",
			CacheEnable:                   true,
			CacheExpired:                  "1d",
			VideoPreviewEnable:            true,
			VideoPreviewContainers:        "mp4,mkv",
			VideoPreviewIgnoreTemplateIds: "LD,SD",
			PathEmby2Openlist:             "/movie:/电影\n/music:/音乐\n/show:/综艺\n/series:/电视剧\n/sport:/运动\n/animation:/动漫",
			LogDisableColor:               true,
			NotifyEnable:                  false,
			NotifyUrl:                     "",
			NotifyMethod:                  "POST",
			NotifyContentType:             "application/json",
			NotifyTitleKey:                "title",
			NotifyContentKey:              "text",
		}).Error
	}
	var m map[string]any
	_ = yaml.Unmarshal(b, &m)
	emby := getMap(m, "emby")
	cache := getMap(m, "cache")
	vp := getMap(m, "video-preview")
	path := getMap(m, "path")
	openlist := getMap(m, "openlist")
	ltg := getMap(openlist, "local-tree-gen")
	strm := getMap(emby, "strm")
	ssl := getMap(m, "ssl")
	g := GlobalConfig{
		EpisodesUnplayPrior:           boolVal(emby, "episodes-unplay-prior", true),
		ResortRandomItems:             boolVal(emby, "resort-random-items", true),
		ProxyErrorStrategy:            strVal(emby, "proxy-error-strategy", "origin"),
		ImagesQuality:                 intVal(emby, "images-quality", 100),
		DownloadStrategy:              strVal(emby, "download-strategy", "403"),
		CacheEnable:                   boolVal(cache, "enable", true),
		CacheExpired:                  strVal(cache, "expired", "1d"),
		VideoPreviewEnable:            boolVal(vp, "enable", true),
		VideoPreviewContainers:        strings.Join(sliceStr(vp, "containers"), ","),
		VideoPreviewIgnoreTemplateIds: strings.Join(sliceStr(vp, "ignore-template-ids"), ","),
		PathEmby2Openlist:             strings.Join(sliceStr(path, "emby2openlist"), "\n"),
		LogDisableColor:               boolVal(getMap(m, "log"), "disable-color", true),
		StrmPathMap:                   strings.Join(sliceStr(strm, "path-map"), "\n"),
		CacheWhiteList:                strings.Join(sliceStr(cache, "whitelist"), "\n"),
		LTGEnable:                     boolVal(ltg, "enable", false),
		LTGFFmpegEnable:               boolVal(ltg, "ffmpeg-enable", false),
		LTGVirtualContainers:          strVal(ltg, "virtual-containers", "mp4,mkv"),
		LTGStrmContainers:             strVal(ltg, "strm-containers", "ts"),
		LTGMusicContainers:            strVal(ltg, "music-containers", "mp3,flac"),
		LTGAutoRemoveMaxCount:         intVal(ltg, "auto-remove-max-count", 6000),
		LTGRefreshInterval:            intVal(ltg, "refresh-interval", 10),
		LTGScanPrefixes:               strings.Join(sliceStr(ltg, "scan-prefixes"), "\n"),
		LTGIgnoreContainers:           strVal(ltg, "ignore-containers", "jpg,jpeg,png,txt,nfo,md"),
		LTGThreads:                    intVal(ltg, "threads", 8),
		SslEnable:                     boolVal(ssl, "enable", false),
		SslSinglePort:                 boolVal(ssl, "single-port", false),
		SslKey:                        strVal(ssl, "key", ""),
		SslCrt:                        strVal(ssl, "crt", ""),
		NotifyEnable:                  false,
		NotifyUrl:                     "",
		NotifyMethod:                  "POST",
		NotifyContentType:             "application/json",
		NotifyTitleKey:                "title",
		NotifyContentKey:              "text",
	}
	return DB.Create(&g).Error
}

func GetGlobalConfig() (GlobalConfig, error) {
	var g GlobalConfig
	err := DB.First(&g).Error
	return g, err
}

func UpdateGlobalConfig(g *GlobalConfig) error {
	return DB.Save(g).Error
}

func getMap(m map[string]any, k string) map[string]any {
	v, _ := m[k]
	mv, _ := v.(map[string]any)
	if mv == nil {
		mv = map[string]any{}
	}
	return mv
}

func strVal(m map[string]any, k string, def string) string {
	if v, ok := m[k]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func boolVal(m map[string]any, k string, def bool) bool {
	if v, ok := m[k]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

func intVal(m map[string]any, k string, def int) int {
	if v, ok := m[k]; ok {
		switch vv := v.(type) {
		case int:
			return vv
		case int64:
			return int(vv)
		case float64:
			return int(vv)
		}
	}
	return def
}

func sliceStr(m map[string]any, k string) []string {
	var out []string
	if v, ok := m[k]; ok {
		switch vv := v.(type) {
		case []any:
			for _, e := range vv {
				if s, ok := e.(string); ok {
					out = append(out, s)
				}
			}
		case []string:
			out = append(out, vv...)
		}
	}
	return out
}
