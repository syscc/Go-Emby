# Changelog

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
