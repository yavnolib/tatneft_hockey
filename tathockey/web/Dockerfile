# Используем официальный образ Go для сборки
FROM golang:1.22.6-alpine AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY back .

# Сборка бинарного файла
RUN go build -o /auth_service cmd/main.go

# Используем минимальный образ для запуска
FROM alpine:latest

# Копируем бинарный файл из этапа сборки
COPY --from=builder /auth_service /auth_service

# Устанавливаем рабочую директорию
WORKDIR /

# Указываем порт, на котором работает сервис
EXPOSE 8000

# Команда запуска контейнера
CMD ["/auth_service"]