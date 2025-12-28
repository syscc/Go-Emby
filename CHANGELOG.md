# Changelog

## v0.0.5
- 修复：WebUI 日志时间解析按本地时区（ParseInLocation + time.Local），与容器 TZ 保持一致
- CI：GitHub Actions 赋予 contents: write 权限，推送 v* 标签自动创建 Release，避免 403

## v0.0.4
- WebUI 全局配置保存后，自动重启所有内核实例（manager.Restart），配置立即生效，无需重启容器
- 更新 CurrentVersion 为 v0.0.4
- 说明：在 Manager 模式下，子模块读取的是各实例专属 config.yml；如需生效应通过 WebUI 的服务器/全局配置界面修改。直接修改根目录 config.yml 仅用于 -kernel-only 模式
