FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rbac-service ./cmd/rbac

FROM alpine:latest

WORKDIR /app

# Kopiowanie skompilowanego pliku z poprzedniego etapu
COPY --from=builder /app/rbac-service /app/

# Port na którym działa aplikacja
EXPOSE 4000

# Komenda uruchamiająca aplikację
CMD ["./rbac-service"]