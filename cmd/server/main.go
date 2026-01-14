package main

import (
	"bufio"
	"chat-server-go/internal/chat"
	"chat-server-go/internal/config"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	cfg := config.LoadConfig()
	addres := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	chat.GenerateTitle("Chat Server in Go", true)

	listener, err := net.Listen("tcp", addres)

	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	defer listener.Close()

	fmt.Printf("Servidor Corriendo en %s\n", addres)

	// Incializar el manejador de usuarios
	hub := chat.NewHub(cfg.MaxConnections)
	go hub.Run()

	// Go routine para monitorear la consola (Input 'exit')
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Ingrese 'exit' para apagar el servidor.")
		for scanner.Scan() {
			if strings.ToLower(strings.TrimSpace(scanner.Text())) == "exit" {
				fmt.Println("Cerrando el servidor...")
				os.Exit(0) // En Go, os.Exit(0) cierra de forma segura
			}
		}
	}()

	// Bucle principal de aceptacion de clientes
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error aceptando la conexion: %v", err)
			continue
		}

		go chat.HandlerConnection(conn, hub)

	}
}
