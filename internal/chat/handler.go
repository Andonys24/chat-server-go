package chat

import (
	"log"
	"net"
	"strings"
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
		case CmdUsers:
			// Crear canal solo para estas respuestas
			responseChan := make(chan string)

			// Pasar hub al canal
			hub.UserListRequest <- responseChan

			// Esperar respuesta (bloqueo momentaneo seguro)
			userList := <-responseChan

			// Enviar respuesta al cliente
			transport.Send(RespInfo, "USERS|"+userList)
		case CmdMessage:
			// El 'content' viene como "Destinatario|Mensaje"
			parts := strings.SplitN(content, "|", 2)
			if len(parts) < 2 {
				transport.Send(RespInfo, InfoTypeError+"|Formato incorrecto. Usa: MESSAGE|usuario|mensaje")
				continue
			}

			hub.PrivateMsg <- &PrivateMessageRequest{
				From:    currentNickname,
				To:      parts[0],
				Content: parts[1],
			}
		case CmdCLeanConsole:
			transport.Send(RespOkClean, "Capa de consola limpia")
		default:
			log.Printf("[ADVERTENCIA] Cliente [%s] envió comando desconocido: %s", currentNickname, header)
		}
	}
}
