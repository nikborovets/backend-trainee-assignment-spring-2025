# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/service/main.go

# Оптимизированный образ для тестов 
FROM golang:1.24-alpine AS tester
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Копируем весь контекст, чтобы избежать проблем с отдельными директориями
COPY . .
# Установка необходимых зависимостей для тестов
RUN apk add --no-cache gcc musl-dev curl
# По умолчанию команда запускает все тесты
CMD ["go", "test", "-v", "./test/..."]

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY .env /app/.env
EXPOSE 8080 3000 9000
CMD ["/app/app"]