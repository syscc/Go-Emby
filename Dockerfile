# 第一阶段：构建阶段
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# 声明构建参数
ARG TARGETOS
ARG TARGETARCH

# 设置工作目录
WORKDIR /app

# 设置代理
RUN go env -w GOPROXY=https://goproxy.cn

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY cmd cmd
COPY internal internal
COPY main.go main.go

# 编译源码成静态链接的二进制文件
# 打印构建环境信息，用于调试
RUN echo "Building for $TARGETOS/$TARGETARCH"
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -tags netgo -ldflags="-s -w -extldflags '-static' -X main.ginMode=release" -o main .

# 第二阶段：运行阶段
# 使用 Alpine 以支持多架构 (linux/386, linux/arm/v7 等)
FROM alpine:latest

# 设置时区和证书
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译后的二进制文件到 /usr/bin，防止被 /app 挂载覆盖
COPY --from=builder /app/main /usr/bin/go-emby

# 赋予执行权限
RUN chmod +x /usr/bin/go-emby

# 暴露端口
EXPOSE 8090

# 运行应用程序
CMD ["/usr/bin/go-emby", "-dr", "/app"]
