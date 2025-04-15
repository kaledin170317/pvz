
FROM golang:1.24

# Переходим в нужную директорию внутри контейнера
WORKDIR /go/pvz

# Копируем всё содержимое проекта (относительно Docker context!)
COPY ../ ./

# Собираем приложение
RUN go build -o /pvz_server ./cmd/pvzapp/main.go \
    && go clean -cache -modcache

EXPOSE 1488

CMD ["/pvz_server"]
