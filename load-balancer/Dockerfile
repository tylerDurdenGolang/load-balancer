# Этап сборки
FROM golang:1.19 AS builder

# Создадим рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем модули
RUN go mod download

# Копируем оставшиеся файлы проекта
COPY . .

# Сборка бинарного файла из пакета cmd/lb
RUN go build -o /load-balancer ./cmd/lb

# Финальный образ
FROM alpine:3.17

# Создадим рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарник из предыдущего этапа
COPY --from=builder /load-balancer /app/load-balancer

# Запускаем бинарник по умолчанию
CMD ["/app/load-balancer"]
