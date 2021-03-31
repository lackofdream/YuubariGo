FROM golang:alpine AS builder

RUN echo https://mirrors.sjtug.sjtu.edu.cn/alpine/v3.12/main > /etc/apk/repositories && apk update && apk add --no-cache git

WORKDIR $GOPATH/src/yuubari_go
COPY . .

RUN go build -mod=vendor -ldflags "-w -s" -o /go/bin/yuubari_go yuubari_go/cli

FROM alpine:latest

RUN echo https://mirrors.sjtug.sjtu.edu.cn/alpine/v3.12/main > /etc/apk/repositories && apk add --no-cache tzdata
ENV TZ Asia/Shanghai

COPY --from=builder /go/bin/yuubari_go /go/bin/yuubari_go

ENTRYPOINT ["/go/bin/yuubari_go"]

CMD ["-debug", "-interval", "2", "-retry", "10", "-kcp", "http://127.0.0.1:8081"]
