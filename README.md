<div align="center">
  <img height="150px" src="./internal/webui/static/logo.png"></img>
</div>

<h1 align="center">Go-Emby</h1>

<div align="center">
  <a href="https://github.com/syscc/Go-Emby/tree/v0.0.1"><img src="https://img.shields.io/github/v/tag/syscc/Go-Emby"></img></a>
  <a href="https://github.com/syscc/Go-Emby/releases/latest"><img src="https://img.shields.io/github/downloads/syscc/Go-Emby/total"></img></a>
  <img src="https://img.shields.io/github/stars/syscc/Go-Emby"></img>
  <img src="https://img.shields.io/github/license/syscc/Go-Emby"></img>
</div>

<div align="center">
  Go 语言编写的 Emby + OpenList 网盘直链反向代理服务，深度适配阿里云盘转码播放。
</div>

## ✨ 更新亮点：WebUI 管理后台

本项目已全面重构，引入了全新的 WebUI 管理后台，告别繁琐的手动配置文件编辑：

- **一站式配置与管理**
  - 可在网页端轻松添加和编辑多个媒体服务器。
  - 支持可视化的配置项：名称、端口、Emby Host/Token、OpenList Host/Token、挂载路径、内部重定向、直链缓存时间等。
  - **独立进程隔离**：每个媒体服务器配置保存后，会生成独立的配置文件（`/app/servers/<Name>/config.yml`）并以独立的子进程运行核心内核，互不干扰。
  - **即时生效**：在 WebUI 中激活/修改媒体服务器后，会自动重启相应的子进程。

- **日志查看与过滤**git status
  - 支持在 WebUI 中实时查看系统日志。
  - 支持按服务器筛选日志，并提供多种过滤类别（播放、错误、直链、重定向、字幕等）。

- **统一的持久化存储**
  - 所有数据（数据库、日志、各服务器配置、自定义脚本/样式）统一存储在 `/app` 目录，迁移和备份更加方便。

## 📖 小白必看

**网盘直链反向代理原理**:

正常情况下，Emby 通过磁盘挂载读取网盘资源，流量经过 Emby 服务器中转，速度受限于服务器带宽和性能。

使用网盘直链反向代理后：
1. **Emby Api 请求**：仍由反代服务器转发给 Emby 源服务器。
2. **视频播放请求**：反代服务器向 OpenList 获取网盘直链，并直接重定向返回给客户端。
3. **直连播放**：客户端拿到直链后，直接从网盘下载数据进行解码播放，**不再经过 Emby 服务器**。

**优势**：
- 播放速度取决于你的网速和网盘带宽（通常能拉满）。
- 减轻 Emby 服务器的 CPU 和带宽压力。
- 支持 4K 原画流畅播放（取决于客户端解码能力）。

## 🚀 功能特性

- **OpenList 网盘原画直链播放**
- **Strm 直链播放**
- **[OpenList 本地目录树生成](#-使用说明-openlist-本地目录树生成)**：自动扫描 OpenList 生成本地 strm/虚拟文件。
- **[自定义注入 js/css](#-使用说明-自定义注入-web-jscss)**：支持 Web 端自定义脚本和样式。
- **阿里云盘转码直链播放**：
  - 不消耗三方流量包（非会员具体情况需自测）。
  - 兼容性好，支持 Web、AndroidTV 等多种客户端。
  - *局限*：多音轨只能播放默认音轨，内封字幕丢失（支持外挂/转码字幕）。
- **Websocket 代理**
- **客户端防转码（转容器）**
- **缓存中间件**：直链缓存（默认 10 分钟）、字幕缓存（30 天）、API 缓存。

## ✅ 已测试并支持的客户端

| 名称 | 原画支持 | 阿里转码支持 | 说明 |
| :--- | :---: | :---: | :--- |
| [`Gemby`](https://github.com/AmbitiousJun/gemby) | ✅ | ✅ | 推荐使用 |
| `Emby Web` | ✅ | ✅ | 转码字幕有概率挂载不上 |
| `Emby for Android` | ✅ | ✅ | |
| `Emby for AndroidTV` | ✅ | ✅ | 遥控器调进度可能会触发频繁请求限制 |
| `Fileball` | ✅ | ✅ | |
| `Infuse` | ✅ | ❌ | 建议设置缓存方式为`不缓存` |
| `VidHub` | ✅ | ✅ | |
| `Emby for Kodi Next Gen` | ✅ | ✅ | 需开启插件设置：播放/视频转码/prores |

## 🛠 前置环境准备

1. **Emby 服务器** & **OpenList 服务器**。
2. **Docker** 环境。
3. **本地挂载**：通过 cd2 或 rclone 将网盘挂载到本地，供 Emby 读取（推荐使用 cd2 连接阿里云盘，并将缓存设为极小值如 1MB）。
4. **路径对应**：确保 Emby 媒体库路径与 OpenList 挂载路径能对应上（可通过配置中的路径映射 `path.emby2openlist` 解决）。

## 🐳 使用 DockerCompose 部署

**不再需要手动复制和编辑 config.yml 文件，一切配置均在 WebUI 中完成。**

### 1. 创建目录与文件

在服务器上创建一个目录（如 `go-emby`），并在其中创建 `docker-compose.yml` 文件：

```yaml
version: "3.8"
services:
  go-emby:
    image: syscc/go-emby:latest
    container_name: go-emby
    restart: always
    environment:
      - TZ=Asia/Shanghai
    volumes:
      # 数据持久化目录（包含数据库、日志、配置、自定义脚本等）
      - ./app:/app
      # 如果使用了 OpenList 本地目录树生成功能，需要映射对应目录
      # - ./openlist-local-tree:/app/openlist-local-tree
    network_mode: host
    # 如果不支持 host 模式（如 Docker Desktop），请使用端口映射：
    # ports:
    #   - 8090:8090  # WebUI 管理端口
    #   - 8095:8095  # HTTP 内核默认端口（可在 WebUI 修改）
    #   - 8094:8094  # HTTPS 内核默认端口（可在 WebUI 修改）
```

### 2. 启动服务

```shell
docker-compose up -d
```

### 3. 初始化配置

1. 浏览器访问 `http://ip:8090` 进入 WebUI 管理后台。
2. 首次访问需创建管理员账号。
3. 进入后台后，在“全局配置”中设置 OpenList 本地目录树等全局参数（可选）。
4. 点击“添加服务器”，配置你的 Emby 和 OpenList 信息。
   - **名称**：任意起名（如 `MyEmby`）。
   - **Go Emby 端口**：设置该实例的监听端口（默认 8095）。
   - **Emby Host/Token**：你的 Emby 服务器地址和 API Key。
   - **OpenList Host/Token**：你的 OpenList 地址和密码。
   - **挂载路径**：配置路径映射规则。
5. 保存后，服务将自动启动。

## 📝 使用说明

### 🔐 SSL 配置

在 WebUI 的 **全局配置** 中设置 SSL：
1. 将证书 (`.crt`) 和私钥 (`.key`) 文件放入映射的 `./app/ssl` 目录中（如果目录不存在请手动创建）。
2. 在 WebUI 全局配置中填入证书和私钥的**绝对路径**（容器内路径，例如 `/app/ssl/fullchain.crt` 和 `/app/ssl/privkey.key`）。
3. 保存并重启相关的媒体服务器实例。

### 🎨 自定义注入 Web JS/CSS

将自定义的 `.js` 或 `.css` 文件放入映射的 `./app/custom-js` 或 `./app/custom-css` 目录中，重启服务即可自动注入到 Emby Web 端。

| 描述 | 推荐脚本/样式 |
| :--- | :--- |
| 外部播放器按钮 | [ExternalPlayers.js](https://emby-external-url.7o7o.cc/embyWebAddExternalUrl/embyLaunchPotplayer.js) |
| 首页轮播图 | [emby-swiper.js](https://github.com/AmbitiousJun/emby-css-js/raw/refs/heads/main/custom-js/emby-swiper.js) |
| 节目界面美化 | [show-display.css](https://github.com/AmbitiousJun/emby-css-js/raw/refs/heads/main/custom-css/show-display.css) |

### 🌲 OpenList 本地目录树生成

在 WebUI 的 **全局配置** -> **本地目录树生成** 中开启并配置：
- **启用**：开启功能。
- **FFmpeg 启用**：开启后可提取视频真实时长和音乐元数据（需下载 ffmpeg 环境）。
- **容器格式**：配置需要处理的视频/音频后缀（如 `mp4,mkv`）。
- **刷新间隔**：设置自动扫描的间隔时间。

配置完成后，需确保 `docker-compose.yml` 中映射了 `/app/openlist-local-tree` 目录，以便将生成的 strm/虚拟文件保存到宿主机。

## ☕️ 请我喝杯咖啡

<img height="300px" src="./assets/pay.jpg"></img>

## 🌟 Star History

[![Stargazers over time](https://starchart.cc/syscc/Go-Emby.svg?variant=adaptive)](https://starchart.cc/syscc/Go-Emby)

