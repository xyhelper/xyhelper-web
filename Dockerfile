# 使用 golang 1.16 作为基础镜像
FROM golang:1.18 AS build

# 将当前工作目录设置为 /app
WORKDIR /app

# 将源代码复制到容器内的 /app 目录
COPY . .

# 配置goproxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 编译应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# 使用 alpine 作为基础镜像
FROM alpine:latest

# 将时区设置为 Asia/Shanghai
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 将应用程序复制到容器内
COPY --from=build /app/app /usr/local/bin/app

# 设置环境变量
ENV PORT=8080

# 暴露端口
EXPOSE $PORT

# 启动应用程序
CMD ["app"]
