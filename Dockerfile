# 构建阶段
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用，确保静态链接，明确指定 arm64 架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o semrel-gitlab

# 运行阶段
FROM ubuntu:22.04

# 安装时区数据
RUN apt-get update && apt-get install -y tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/lib/apt/lists/*

# 从构建阶段复制二进制文件
COPY --from=builder /app/semrel-gitlab /usr/local/bin/semrel-gitlab

# 设置入口点
ENTRYPOINT ["/usr/local/bin/semrel-gitlab"]