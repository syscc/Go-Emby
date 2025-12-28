# Changelog

## v0.0.4
- WebUI 全局配置保存后，自动重启所有内核实例（manager.Restart），配置立即生效，无需重启容器
- 更新 CurrentVersion 为 v0.0.4
- 说明：在 Manager 模式下，子模块读取的是各实例专属 config.yml；如需生效应通过 WebUI 的服务器/全局配置界面修改。直接修改根目录 config.yml 仅用于 -kernel-only 模式

