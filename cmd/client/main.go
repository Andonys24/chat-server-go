package main

import (
	"chat-server-go/internal/config"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("Conectando al servidor: %s:%d", cfg.Host, cfg.Port)
}
