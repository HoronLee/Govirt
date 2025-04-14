FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN CGO_ENABLED=0 GOOS=linux go build -o govirt

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/govirt /app/govirt

RUN chmod +x /app/govirt

CMD ["/app/govirt"]