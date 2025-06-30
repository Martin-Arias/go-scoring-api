# Build the Go application
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# runtime image
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]