package chat

import (
	"fmt"
	"net"
)

func HandlerConnection(conn net.Conn, hub *Hub) {
	transport := NewTransport(conn)

	// Fase de Registro
	header, content, err := transport.Receive()
	if err != nil || header != CmdEnter {
		transport.Send(RespErrorEnter, "Identificación requerida")
		conn.Close()
		return
	}

	currentNickname := content

	// Validación rápida antes de molestar al Hub
	if !IsValidNickname(currentNickname) {
		transport.Send(RespErrorEnter, "Nickname inválido (3-12 caracteres, empieza con letra)")
		conn.Close()
		return
	}

	// Pedir registro al Hub
	hub.Register <- &RegisterRequest{
		Nickname:  currentNickname,
		Conn:      conn,
		Transport: transport,
	}

	defer func() {
		hub.Unregister <- conn
	}()

	// BUCLE INFINITO: Mantiene la conexión abierta y escuchando
	for {
		header, content, err = transport.Receive()
		if err != nil {
			// Si el cliente cierra la terminal o falla la red
			hub.Unregister <- conn
			return
		}

		switch header {
		case CmdExit:
			// Si el usuario escribe EXIT, salir
			return
		case CmdAll:
			hub.Broadcast <- &MessageRequest{
				From:    currentNickname,
				Content: content,
			}
		}
	}
}

func (h *Hub) handleRemoveUser(conn net.Conn) {
	var nicknameToRemove string

	for name, user := range h.Clients {
		if user.Conn == conn {
			nicknameToRemove = name
			break
		}
	}

	if nicknameToRemove != "" {
		// Eliminar del mapa
		delete(h.Clients, nicknameToRemove)

		// Cerrar la conexion fisicamente
		conn.Close()

		fmt.Printf("Hub: Usuario [%s] desconectado. Total: %d\n", nicknameToRemove, len(h.Clients))

		// Notificar a los demas INFO|EXIT|Nombre
		h.broadcastInfo(InfoTypeExit, nicknameToRemove, nil)
	}

	conn.Close()
}
