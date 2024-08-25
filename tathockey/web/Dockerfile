FROM golang:1.22.6-alpine AS build

WORKDIR /var/backend

# Копируем go.mod и go.sum для установки зависимостей
COPY cmd cmd
COPY internal internal
COPY tmp tmp
COPY .env .env
COPY api api

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod tidy
RUN go build -o main ./cmd/test/main.go

FROM alpine:edge as prod
RUN apk add bash
COPY --from=build /var/backend/main /app/main
COPY --from=build /var/backend/.env /app/.env

COPY tmp /app/tmp
COPY static /app/static
COPY gifs /app/gifs

WORKDIR /app
EXPOSE 8080
ENTRYPOINT ./main