# -------- Этап сборки --------
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем go.mod и go.sum отдельно для кеша зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Сборка бинарника без CGO, для запуска в Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /load-balancer ./cmd/lb

# -------- Минималистичный финальный образ --------
FROM alpine:3.17

WORKDIR /app

# Копируем бинарник из билдера
COPY --from=builder /load-balancer ./load-balancer

# Устанавливаем точку входа
CMD ["./load-balancer"]
