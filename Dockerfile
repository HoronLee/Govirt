FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN CGO_ENABLED=0 GOOS=linux go build -o GoHub

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/GoHub /app/GoHub

RUN chmod +x /app/GoHub

CMD ["/app/GoHub"]