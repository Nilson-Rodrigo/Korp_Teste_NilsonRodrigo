FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código-fonte
COPY . .

# Build otimizado
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app .

# Stage final
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk add --no-cache ca-certificates

# Copiar binário do builder
COPY --from=builder /app/app .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Expor porta
EXPOSE 8081

# Comando para executar
CMD ["./app"]
