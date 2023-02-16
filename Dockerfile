FROM golang:1.19.5

ENV GO111MODULE=on \
    GOPROXY=goproxy.io \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /project/go-docker

# 复制go.mod，go.sum并且下载依赖
COPY go.* ./
RUN go mod download

# 复制项目内的所有内容并构建
COPY . .
RUN go build -o /project/go-docker/build/myapp .

EXPOSE 38088
ENTRYPOINT [ "/project/go-docker/build/myapp" ]