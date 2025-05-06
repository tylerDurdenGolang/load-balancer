# Этап сборки
FROM golang:1.23 AS builder
# Создадим рабочую директорию внутри контейнера
WORKDIR /app
# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./
# Загружаем модули
RUN go mod download
# Копируем оставшиеся файлы проекта
COPY . .
# Сборка бинарного файла из пакета cmd/lb
RUN go build -o /template-api ./cmd
# Финальный образ
FROM alpine:3.17
# Создадим рабочую директорию
WORKDIR /app
# Копируем скомпилированный бинарник из предыдущего этапа
COPY --from=builder /template-apir /app/template-api
# Запускаем бинарник по умолчанию
CMD ["/app/template-api"]
