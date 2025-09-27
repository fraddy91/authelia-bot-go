# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o bot main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
COPY config/ config/
COPY notifications/ notifications/
CMD ["./bot"]