# Используем официальный образ Golang для создания артефакта сборки.
# Это стадия сборки
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go mod и sum файлы
COPY go.mod go.sum ./

# Загружаем все зависимости. Зависимости будут кэшироваться, если go.mod и go.sum файлы не изменялись
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Сборка Go приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/sso

# Используем минимальный образ для запуска нашего приложения
FROM scratch

# Копируем скомпилированное приложение из стадии сборки
COPY --from=builder /main /main
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/.env /.env

# Устанавливаем рабочую директорию
WORKDIR /

# Устанавливаем переменные окружения
ENV CONFIG_PATH=/config/local.yaml
#ENV MIGRATIONS_PATH=/migrations

# Сообщаем Docker, что контейнер слушает этот порт
EXPOSE 5001

# Команда для запуска приложения
CMD ["/main"]
