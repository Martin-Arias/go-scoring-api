# Builder con Go oficial
FROM golang:1.24-alpine AS builder

# Instalar herramientas necesarias en etapa de build
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar binario estático
RUN CGO_ENABLED=0 go build -o main ./cmd

# Final image: Alpine
FROM alpine:3.20

WORKDIR /app

# Copiar binario desde build
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

ENTRYPOINT ["./main"]