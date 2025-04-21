# syntax=docker/dockerfile:1

FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/service/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY .env /app/.env
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/app/app"] 