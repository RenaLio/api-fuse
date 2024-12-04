# 第一阶段：Builder
FROM golang:1.23.3 AS builder

# 设置工作目录
WORKDIR /app

# 复制Go模块文件
COPY go.mod go.sum ./

# 设置镜像源
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译Go应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mainApp ./cmd/api_app/


# 第二阶段：Alpha镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从Builder阶段复制编译后的二进制文件
COPY --from=builder /app/mainApp .

# 设置可执行权限
RUN chmod +x ./mainApp

# 设置环境变量（如果需要）
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release
# 暴露端口（如果需要）
EXPOSE 8080

# 运行应用程序
CMD ["./mainApp"]

#docker build -t chat-app:0.0.5 .

#docker save chat-app:0.0.5 -o .\release\chat_v0.0.5.tar
# docker load -i ./chat_v0.0.5.tar
#docker run -d --name chat-dev -p 19999:8080 -v /e/docker/chat_dev:/app/config --restart=always chat-app:0.0.4

