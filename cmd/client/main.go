package main

import (
	"bufio"
	"chat-server-go/internal/chat"
	"chat-server-go/internal/config"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	cfg := config.LoadConfig()
	addres := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))

	// Conexion al servidor
	conn, err := net.Dial("tcp", addres)

	if err != nil {
		fmt.Println("Error al conectar al servidor:", err)
		return
	}

	defer conn.Close()

	chat.GenerateTitle("Chat Server By Andonys24", true)

	fmt.Println("Conectado al servidor de Chat.")

	// Lanzar una Goroutine q ESCUCHA al servidor
	go listenToServer(conn)

	// El Hilo principal se encarga de Leer el teclado del cliente
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if strings.TrimSpace(input) == "" {
			continue
		}

		// Enviar data al servidor
		_, err := fmt.Fprintf(conn, "%s\n", input)

		if err != nil {
			fmt.Println("Error enviando mensaje:", err)
			break
		}

		if input == chat.CmdExit {
			break
		}
	}
}

// ListenToServer se encarga de imprimir lo que llega del servidor
func listenToServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\n\033[31m[!] Conexión perdida con el servidor.\033[0m")
			os.Exit(0)
		}

		msg := strings.TrimSpace(message)
		parts := strings.Split(msg, "|")
		header := parts[0]

		// Limpiamos la línea actual para que el mensaje no se mezcle con el >>>
		fmt.Print("\r\033[K")

		switch header {
		case chat.RespOkEnter:
			fmt.Printf("\033[32m[SISTEMA]: %s\033[0m\n", parts[1])

		case chat.RespInfo:
			// parts[1] es ENTER/EXIT, parts[2] es el nombre
			switch parts[1] {
			case chat.CmdEnter:
				fmt.Printf("\033[34m(i) %s se ha unido al chat.\033[0m\n", parts[2])
			case chat.CmdExit:
				fmt.Printf("\033[33m(i) %s ha salido del chat.\033[0m\n", parts[2])
			case chat.CmdUsers:
				fmt.Printf("\033[36m[USUARIOS]: %s\033[0m\n", parts[2])
			default:
				fmt.Printf("\033[34m[INFO]: %s\033[0m\n", parts[2])
			}

		case chat.RespMsgFrom: // Mensaje Privado o Global
			// Formato: OF|Remitente|Contenido
			fmt.Printf("\033[35m[Mensaje de %s]:\033[0m %s\n", parts[1], parts[2])

		case chat.RespOkClean:
			// Comando para limpiar pantalla (ANSI escape)
			chat.CleanConsole()
			fmt.Println("--- Consola Limpia ---")

		default:
			fmt.Println(msg)
		}

		// Volvemos a poner el prompt
		fmt.Print(">>> ")
	}
}
