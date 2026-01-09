FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/migrator ./cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bot ./cmd/bot

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/migrator /app/migrator
COPY --from=builder /app/bot /app/bot

COPY internal/storage/postgres/migrations ./internal/storage/postgres/migrations
COPY .env .

LABEL authors="recrusion"

CMD ["/app/bot"]