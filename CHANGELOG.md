# Changelog

## v0.1.0
- 核心修复：解决高并发下的严重内存泄漏问题
  - 修复 HTTP 重定向响应体未关闭导致的连接/内存堆积
  - 缓存写入增加 5MB 阈值保护，防止大文件（如视频流）被无脑读入内存
  - 优化 HTTP Client 连接池配置（MaxIdleConns 提升至 100），减少 TIME_WAIT
- 缓存策略：默认缓存策略收紧，仅缓存直链重定向（ResourceStream/ResourceOriginal），不再缓存 PlaybackInfo/字幕/下载等非直链接口
- 数据库：数据库文件重命名为 `Go-Emby.db`（原 `Emby-Go.db`），启动时自动迁移
- 版本：CurrentVersion 提升为 v0.1.0

## v0.0.9
- STRM 直链缓存：新增“忽略缓存域名”，支持多域名/通配符/关键字，仅匹配域名部分（如 *.115cdn.*、115）
- WebUI：服务器弹窗新增“忽略缓存域名”文本域，支持多行输入；保存后写入各实例配置
- Manager：写入 config.yml 增加 emby.dl-cache-ignore（支持逗号/分号/换行分隔）
- Kernel：实现 isCacheIgnored（通配符与关键字匹配），命中后显式设置 Expired=-1，避免二次播放出现“直链缓存命中”
- 版本：CurrentVersion 提升为 v0.0.9

## v0.0.8
- 健康检查：每分钟检查子进程是否返回 50x；连续复检 2 次（间隔 2 秒）仍异常则重启；重启后仍异常发送通知；全流程增加日志提示
- 代理异常：反代上游不可达时返回 502（Bad Gateway），便于健康检查识别异常
- 通知配置：支持新增多条通知并列表管理，新增“通知名称”字段；测试消息校验启用与地址；支持 application/json 与 x-www-form-urlencoded
- WebUI：通知配置改为“新增弹框”交互，保存后列表显示并可编辑/删除

## v0.0.7
- WebUI：品牌名统一为 “Go-Emby”（标题、侧边栏、文案）
- STRM：内部重定向优先使用 HEAD（3 次，1s 间隔），失败回退 GET（2 次，1s 间隔），兜底保留原始链接；不影响缓存逻辑
- CI：发布附件不再包含 README/LICENSE，仅保留版本化二进制与校验文件

## v0.0.6
- 修复：转发到 Emby 时保留上游 X-Forwarded-For 链并追加当前客户端 IP；X-Real-IP 取链首，恢复通知/播放日志显示公网或 CDN IP

## v0.0.5
- 修复：WebUI 日志时间解析按本地时区（ParseInLocation + time.Local），与容器 TZ 保持一致
- CI：GitHub Actions 赋予 contents: write 权限，推送 v* 标签自动创建 Release，避免 403
- WebUI：系统日志移除“时间”列，仅显示“服务器名称/级别/消息”；消息内时间保留且正确
- WebUI：清空日志后设置截止时间，后续刷新不再回显清空前的历史日志
- Web/Kernel：默认运行模式改为 Release；Gin 设置不信任任何代理，移除 “[GIN-debug]” 与代理警告日志

## v0.0.4
- WebUI 全局配置保存后，自动重启所有内核实例（manager.Restart），配置立即生效，无需重启容器
- 更新 CurrentVersion 为 v0.0.4
- 说明：在 Manager 模式下，子模块读取的是各实例专属 config.yml；如需生效应通过 WebUI 的服务器/全局配置界面修改。直接修改根目录 config.yml 仅用于 -kernel-only 模式
