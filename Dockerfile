FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o url-shortener ./cmd/url-shortener/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y libc6 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

COPY --from=builder /app/url-shortener /url-shortener

COPY config /app/config
COPY storage/storage.db /app/storage/storage.db

ENV PORT=8080
#CONFIG_PATH=/app/config/prod.yaml
#CONFIG_PATH=/app/config/local.yaml
#HTTP_SERVER_PASSWORD=123

ENTRYPOINT ["/url-shortener"]
