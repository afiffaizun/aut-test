# Stage 1: Builder - build locally and copy binary
FROM golang:1.24 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

ENV GOTOOLCHAIN=auto

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/server

# Stage 2: Runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/auth-service .

RUN useradd -m -u 1000 appuser && chown appuser:appuser /app

USER appuser

EXPOSE 8080

ENV GIN_MODE=release

CMD ["./auth-service"]