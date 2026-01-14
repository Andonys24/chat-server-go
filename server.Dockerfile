# Compilacion
FROM golang:1.25.5-alpine AS builder

# Instalacion de dependencias necesarias
RUN apk add --no-cache git

# Directorio de Trabajo
WORKDIR /app

# Copiar archivos de dependencias primero para aprovechar el cache
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar resto del codigo
COPY . .

# Compilar el binario de forma estatica
RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-server ./cmd/server/main.go

# Imagen Final (ligera)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar solo el binario desde el builder
COPY --from=builder /chat-server .

EXPOSE 8080

CMD ["./chat-server"]