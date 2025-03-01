FROM golang:1.20-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./... -v

RUN go build -o order-service .

FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/order-service /app/

CMD ["/app/order-service"]