#运行环境
FROM golang:1.16.15-alpine3.15

#设置环境变量
ENV GO111MODULE=on \
    GCO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#切换到工作目录
WORKDIR /build

#复制代码文件到容器中
COPY . .

#编译代码为app
RUN go build -o app

#移动到存储二进制的目录
WORKDIR /dist

#复制二进制文件
RUN cp /build/app .

#声明端口
EXPOSE 8888

#运行二进制文件
CMD ["dist/app"]

