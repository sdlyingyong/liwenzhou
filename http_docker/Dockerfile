#运行环境
FROM golang:1.16.15-alpine3.15 AS builder

#设置环境变量
ENV GO111MODULE=auto \
    GCO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#切换到工作目录
WORKDIR /build

#复制代码文件到容器中
COPY . .

#编译代码为app
RUN go build -o app .

###################
# 接下来创建一个小镜像
###################
FROM scratch

#从编译镜像中拷贝出来编译好的程序
COPY --from=builder /build/app /

#运行程序
ENTRYPOINT ["/app"]

#整个运行环境+二进制程序 只需要6.18M
#原本需要300mb