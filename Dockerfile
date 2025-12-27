# 第一阶段：构建阶段
# 使用 Debian 基础镜像，构建环境更稳定
FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

# 声明构建参数
ARG TARGETOS
ARG TARGETARCH

# 设置工作目录
WORKDIR /app

# 设置代理 (为了构建稳定，使用官方代理)
RUN go env -w GOPROXY=https://proxy.golang.org,direct

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

# 验证二进制文件类型
RUN file main || true
RUN ls -lh main

# 第二阶段：运行阶段
# 暂时使用 Ubuntu 以排除 Alpine 兼容性问题 (exec user process caused: no such file or directory 通常是动态链接库缺失导致)
FROM ubuntu:22.04

# 设置时区和证书
ENV TZ=Asia/Shanghai
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译后的二进制文件
COPY --from=builder /app/main /app/main

# 赋予执行权限
RUN chmod +x /app/main

# 暴露端口
EXPOSE 8090

# 运行应用程序
CMD ["/app/main"]
