package chat

import (
	"fmt"
	"net"
)

func HandlerConnection(conn net.Conn, hub *Hub) {
	// Cerrar la conexion al final de la funcion
	defer conn.Close()

	// Usar Transport para enviar info al cliente
	transport := NewTransport(conn)

	fmt.Printf(">>> Nueva conexion detectada desde %s\n", conn.RemoteAddr().String())

	// Enviar un saludo con el protocolo: HEADER|CONTENT
	err := transport.Send("INFO", "Conexion exitosa. El servido te desconectara ahora.")

	if err != nil {
		fmt.Printf("Error al saludar al cliente: %v\n", err)
		return
	}

	fmt.Printf("<<< Cliente %s saludado y desconectado.\n", conn.RemoteAddr().String())
}
