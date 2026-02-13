# 通用 Dockerfile
FROM golang:1.25-alpine AS builder

# 传入构建参数：服务入口文件路径 (如 services/user/api/user.go)
ARG SERVICE_PATH
# 传入构建参数：配置文件路径 (如 services/user/api/etc/user-api.yaml)
ARG CONFIG_FILE

WORKDIR /app

# 设置代理 (国内加速可选)
ENV GOPROXY=https://goproxy.cn,direct

# 缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
# 去除 PATH 中的文件名，只保留目录用于 build (go build -o /server ./services/user/api)
# 这里直接用 .go 文件路径的父目录
RUN dir=$(dirname "${SERVICE_PATH}") && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./$dir

# 运行时镜像
FROM alpine:latest

ARG CONFIG_FILE

WORKDIR /app

# 基础工具
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /server /server
COPY --from=builder /app/${CONFIG_FILE} /etc/config.yaml

CMD ["/server", "-f", "/etc/config.yaml"]
