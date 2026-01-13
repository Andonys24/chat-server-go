package main

import (
	"chat-server-go/internal/chat"
	"chat-server-go/internal/config"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()

	chat.GenerateTitle("Chat Server", true)

	fmt.Printf("Corriendo Servidor: %s:%d", cfg.Host, cfg.Port)
}
