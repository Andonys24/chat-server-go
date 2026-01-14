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
	"time"
)

func main() {
	// Carga de configuración
	cfg := config.LoadConfig()
	addres := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// Título visual
	chat.GenerateTitle("Chat Server in Go", true)

	// Inicializar el Hub (Cerebro)
	hub := chat.NewHub(cfg.MaxConnections)
	go hub.Run()

	// Iniciar el servidor TCP en una Goroutine secundaria
	go func() {
		listener, err := net.Listen("tcp", addres)
		if err != nil {
			log.Fatalf("[ERROR] No se pudo iniciar el servidor: %v", err)
		}
		defer listener.Close()

		log.Printf("[SISTEMA] Servidor iniciado y escuchando en %s", addres)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[ERROR] Fallo al aceptar conexión entrante: %v", err)
				continue
			}
			// Cada cliente en su propia habitación (Goroutine)
			go chat.HandlerConnection(conn, hub)
		}
	}()

	// Bucle de comandos del Administrador (Hilo Principal)
	fmt.Println("Ingrese 'exit' para apagar el servidor.")

	time.Sleep(100 * time.Millisecond)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\rADMIN > ")
		if scanner.Scan() {
			text := strings.ToUpper(strings.TrimSpace(scanner.Text()))

			if text == "" {
				continue
			}

			if text == chat.CmdExit {
				log.Println("[SISTEMA] Apagando el servidor por comando administrativo...")
				os.Exit(0)
			}

			// Si el comando no es EXIT, podrías añadir otros aquí
			log.Printf("[ADVERTENCIA] Comando administrativo no reconocido: %s", text)
		}
	}
}
