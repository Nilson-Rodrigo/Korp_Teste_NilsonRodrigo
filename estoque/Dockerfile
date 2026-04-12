FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código-fonte
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage final
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar binário do builder
COPY --from=builder /app/main .

# Expor porta (será sobrescrita por docker-compose)
EXPOSE 8080

CMD ["./main"]
