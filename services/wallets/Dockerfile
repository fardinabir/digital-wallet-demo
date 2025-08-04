# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Final stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config.docker.yaml .
COPY --from=builder /app/migrations ./migrations/

EXPOSE 8081
CMD ["./main", "server", "--config", "config.docker.yaml"]