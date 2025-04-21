
FROM golang:1.24


WORKDIR /go/pvz
COPY ../ ./

RUN go build -o /pvz_server ./cmd/pvzapp/main.go \
    && go clean -cache -modcache

# REST API порт
EXPOSE 8080

# gRPC порт
EXPOSE 3000

# Prometheus  порт
EXPOSE 9000

CMD ["/pvz_server"]
