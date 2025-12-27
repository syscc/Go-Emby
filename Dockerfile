# 第一阶段：构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装基础工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源码
COPY . .

# 编译源码成静态链接的二进制文件
# -v 打印编译过程，方便排查错误
RUN CGO_ENABLED=0 go build -v -ldflags="-s -w -X main.ginMode=release" -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 设置时区
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译后的二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8095
EXPOSE 8094
EXPOSE 8090

# 运行应用程序
CMD ["./main"]
