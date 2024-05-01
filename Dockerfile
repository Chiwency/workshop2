# 使用官方Go镜像作为构建环境
FROM golang:1.22 as builder

# 设置工作目录
WORKDIR /app

RUN apt-get update && apt-get install -y make

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码
COPY . .

# 使用Makefile来构建master和worker二进制文件
RUN make build-all

