# 使用Alpine Linux作为基础镜像
FROM golang:1.16 as builder

COPY src /src

WORKDIR /src

# 编译并将二进制文件保存到/app/app
RUN go build -o /app/app

# 切换到更小的Alpine Linux镜像
FROM alpine:latest

# 复制构建的二进制文件
COPY --from=builder /app/app /app/app

# 设置容器内的工作目录
WORKDIR /app

CMD ["./app"]

# 暴露端口
EXPOSE 8180
