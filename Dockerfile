# Build stage
FROM golang:1.25.1-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Устанавливаем goose CLI
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Собираем сервер
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/wallet-service

# Final
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/bin/server /app/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /app/migrations

EXPOSE 8080

CMD ["sh", "-c", "goose -dir /app/migrations postgres \"host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} port=${POSTGRES_PORT} sslmode=${POSTGRES_SSLMODE}\" up && /app/server"]
