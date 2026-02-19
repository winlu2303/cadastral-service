FROM golang:1.21-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git gcc musl-dev

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Копируем конфигурационные файлы
COPY config ./config

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]