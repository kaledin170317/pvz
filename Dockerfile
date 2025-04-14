# Этап сборки
FROM golang:1.24 as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o pvz_server ./main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/pvz_server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 1488
ENTRYPOINT ["./pvz_server"]
