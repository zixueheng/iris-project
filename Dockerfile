# 执行构建命令：docker build -t iris-project-app .
# 启动容器：docker run -d --name iris-project-app --network host synm（linux下可用）
# docker run -d --name iris-project-app -p 8080:8080 synm（确保容器能访问mysql和redis）
# 容器内访问：curl -X POST -d '{"username": "admin", "password": "123456"}' -H 'Content-Type: application/json' http://127.0.0.1:8080/adminapi/login
FROM golang:latest as builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/app
# Caching go modules and build the binary.
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iris-project-app .

FROM alpine:latest
WORKDIR /root/goapp
COPY --from=builder /go/src/app/config/ ./config/
COPY --from=builder /go/src/app/iris-project-app .
# 容器向外提供服务的暴露端口
EXPOSE 8080
# 启动服务
ENTRYPOINT ["./iris-project-app"]