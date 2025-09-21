# Build stage
FROM golang:1.25.1-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git make

# 🟡 ИСПРАВЛЕНО: Убрана установка goose из production образа
# Миграции должны выполняться отдельно во время деплоя, а не включаться в runtime образ
# БЫЛО: RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Рабочая директория
WORKDIR /app

# Копируем только файлы модулей сначала (для лучшего кэширования)
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем ВСЕ исходные файлы (исправление ошибки!)
COPY . .  

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/server

# Final stage
FROM alpine:3.20

# Устанавливаем зависимости для runtime
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN adduser -D -g '' appuser

# Переключаемся на непривилегированного пользователя
USER appuser

# Рабочая директория
WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/url-shortener . 

# 🟡 ИСПРАВЛЕНО: Убрано копирование goose из production образа
# БЫЛО: COPY --from=builder /go/bin/goose /usr/local/bin/goose 

# Копируем миграции (только для reference, не для выполнения)
COPY --from=builder /app/migrations ./migrations/  

# Экспонируем порт
EXPOSE 8080

# Команда запуска
CMD ["./url-shortener"]