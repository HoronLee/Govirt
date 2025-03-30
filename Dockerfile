FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN CGO_ENABLED=0 GOOS=linux go build -o GoHub

FROM alpine:latest

LABEL maintainer="horonlee@foxmail.com"

COPY --from=builder /app/GoHub /usr/local/bin/GoHub

RUN chmod +x /usr/local/bin/GoHub

CMD ["/usr/local/bin/GoHub"]