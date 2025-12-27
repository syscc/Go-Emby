# 第一阶段：构建阶段
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

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
# 使用交叉编译，避免在 QEMU 模拟器中编译，极大提高构建速度
# 强制静态链接：-tags netgo -ldflags '-extldflags "-static"'
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -tags netgo -ldflags="-s -w -extldflags '-static' -X main.ginMode=release" -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 设置时区
RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译后的二进制文件
COPY --from=builder /app/main /app/main

# 暴露端口
EXPOSE 8095
EXPOSE 8094
EXPOSE 8090

# 运行应用程序
CMD ["/app/main"]
