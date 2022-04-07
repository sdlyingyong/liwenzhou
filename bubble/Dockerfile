#运行环境
FROM golang:1.16.15-alpine3.15 AS builder

#设置环境变量
#cgo_enabled 控制是否把外部函数打包到生成文件中
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#切换到工作目录
WORKDIR /build

#下载GO类库
COPY go.mod .
COPY go.sum .
RUN go mod tidy

#复制代码文件到容器中
COPY . .

#编译代码为app
RUN go build -o bluebell_app .

###################
# 接下来创建一个小镜像
###################
FROM scratch

COPY ./templates    /templates
COPY ./static       /static
COPY ./conf         /conf

#从编译镜像中拷贝出来编译好的程序
COPY --from=builder /build/bluebell_app /

#运行程序
ENTRYPOINT ["/bluebell_app","conf/config.ini"]

#整个运行环境+二进制程序 只需要6.18M
#原本需要300mb