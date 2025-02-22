FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY order_service .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd/orders/main.go

# Etap finalny
FROM alpine:latest
WORKDIR /
COPY --from=builder /order-service /order-service
COPY .env .
EXPOSE 3000
CMD ["/order-service"]