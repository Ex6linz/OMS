FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1

ENV SKIP_DB=true

COPY order_service/ .

RUN go mod download

RUN go test ./... -v

RUN go build -o order-service ./cmd/orders

FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/order-service /app/

CMD ["/app/order-service"]