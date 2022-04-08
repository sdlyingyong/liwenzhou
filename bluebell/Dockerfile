#单独Dockerfile 版本
##运行环境
#FROM golang:1.16.15-alpine3.15 AS builder
#
##设置环境变量
##cgo_enabled 控制是否把外部函数打包到生成文件中
##使用国内代理下载包
#ENV GO111MODULE=auto \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64    \
#    GOPROXY=https://goproxy.cn,direct
#
##切换到工作目录
#WORKDIR /build
#
##下载GO类库
#COPY go.mod .
#COPY go.sum .
#RUN go mod tidy
#
##复制代码文件到容器中
#COPY . .
#
##编译代码为app
#RUN go build -o bluebell_app .
#
####################
## 接下来创建一个小镜像,用来运行的镜像
####################
#FROM scratch
#
#COPY ./conf         /conf
#
##从编译镜像中拷贝出来编译好的程序
#COPY --from=builder /build/bluebell_app /
#
##运行程序
#ENTRYPOINT ["/bluebell_app"]
#
##整个运行环境+二进制程序 只需要6.18M
##原本需要300mb



#docker-compose 版本
FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod tidy

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o bubble .

###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./static /static
COPY ./conf /conf


# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/bubble /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
        chmod 755 wait-for.sh

# 需要运行的命令
# ENTRYPOINT ["/bubble", "conf/config.ini"]