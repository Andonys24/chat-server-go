FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-client ./cmd/client/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /chat-client .

# Cliente Interactivo, solo se prepara
ENTRYPOINT ["./chat-client"]