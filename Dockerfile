FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY order_service .

RUN go mod download

RUN go test ./... -v

RUN go build -o order-service .

FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/order-service /app/

CMD ["/app/order-service"]