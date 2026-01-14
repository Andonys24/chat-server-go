package chat

import (
	"fmt"
	"net"
)

type Hub struct {
	Clients map[string]*User

	// Configuracion
	MaxConnections int

	// Canales para comunicacion segura
	Register   chan *RegisterRequest
	Unregister chan net.Conn // Usar el ID para desconectar
	Broadcast  chan *MessageRequest
}

// Estructura auxiliar para pedir registro
type RegisterRequest struct {
	Nickname  string
	Conn      net.Conn
	Transport *Transport
}

type MessageRequest struct {
	From    string
	Content string
}

func NewHub(maxConn int) *Hub {
	return &Hub{
		Clients:        make(map[string]*User),
		MaxConnections: maxConn,
		Register:       make(chan *RegisterRequest),
		Unregister:     make(chan net.Conn),
		Broadcast:      make(chan *MessageRequest),
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub de usuarios iniciado...")
	for {
		select {
		case req := <-h.Register:
			h.handleAddUser(req)
		case conn := <-h.Unregister:
			h.handleRemoveUser(conn)
		case msg := <-h.Broadcast:
			h.handleBroadcastMessage(msg)
		}
	}
}

func (h *Hub) handleAddUser(req *RegisterRequest) {
	// Validar espacio (userCount >= MAX_USERS)
	if len(h.Clients) >= h.MaxConnections {
		req.Transport.Send(RespErrorEnter, "Servidor lleno. Intenta mas tarde.")
		req.Conn.Close()
		return
	}

	// Validar si el nombre ya existe
	if _, exists := h.Clients[req.Nickname]; exists {
		req.Transport.Send(RespErrorEnter, "El nombre de usuario ya esta en uso")
		req.Conn.Close()
		return
	}

	// Crear y agregar
	newUser := NewUser(req.Nickname, req.Conn, req.Transport)
	h.Clients[req.Nickname] = newUser

	// Respuesta de exito
	req.Transport.Send(RespOkEnter, "Bienvenido "+req.Nickname)
	fmt.Printf("Hub: Usuario [%s] conectado. Total: %d\n", req.Nickname, len(h.Clients))

	// Notificar a otros
	h.broadcastInfo(InfoTypeEnter, req.Nickname, newUser)
}

func (h *Hub) broadcastInfo(infoType, message string, exclude *User) {
	for _, user := range h.Clients {
		if user != exclude {
			user.Transport.Send(RespInfo, infoType+"|"+message)
		}
	}
}

func (h *Hub) handleBroadcastMessage(msg *MessageRequest) {
	for name, user := range h.Clients {
		if name != msg.From {
			user.Transport.Send(RespMsgFrom, msg.From+"|"+msg.Content)
		}
	}
}
