package chat

import (
	"fmt"
	"net"
	"strings"
)

type Hub struct {
	Clients map[string]*User

	// Configuracion
	MaxConnections int

	// Canales para comunicacion segura
	Register        chan *RegisterRequest
	Unregister      chan net.Conn // Usar el ID para desconectar
	Broadcast       chan *MessageRequest
	UserListRequest chan chan string            // Un canal que recibe canles de strings
	PrivateMsg      chan *PrivateMessageRequest // Canal para procesar peticiones privadas
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

type PrivateMessageRequest struct {
	From    string
	To      string
	Content string
}

func NewHub(maxConn int) *Hub {
	return &Hub{
		Clients:         make(map[string]*User),
		MaxConnections:  maxConn,
		Register:        make(chan *RegisterRequest),
		Unregister:      make(chan net.Conn),
		Broadcast:       make(chan *MessageRequest),
		UserListRequest: make(chan chan string),
		PrivateMsg:      make(chan *PrivateMessageRequest),
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
		case responseChan := <-h.UserListRequest:
			// El Hub genera la lista de forma segura porque está en su propio hilo
			list := h.generateUserList()
			// Envía la respuesta al canal que el Handler le pasó
			responseChan <- list
		case pMsg := <-h.PrivateMsg:
			h.handlePrivateMessage(pMsg)
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

func (h *Hub) handlePrivateMessage(req *PrivateMessageRequest) {
	targetUser, exists := h.Clients[req.To]
	sender, senderExists := h.Clients[req.From]

	if !exists {
		if senderExists {
			sender.Transport.Send(RespInfo, InfoTypeError+"|Usuario '"+req.To+"' no encontrado")
		}
		return
	}

	// Enviar mensaje al destinatario usando el formato OF|Remitente|Mensaje
	targetUser.Transport.Send(RespMsgFrom, req.From+"|"+req.Content)

	// Confirmar al remitente que se envió con éxito (INFO|SUCCESS|...)
	sender.Transport.Send(RespInfo, InfoTypeSuccess+"|Mensaje privado enviado a "+req.To)
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

// Funcion auxiliar para generar string
func (h *Hub) generateUserList() string {
	if len(h.Clients) == 0 {
		return "0"
	}

	var names []string
	for name := range h.Clients {
		names = append(names, name)
	}
	return fmt.Sprintf("%d %s", len(h.Clients), strings.Join(names, " "))
}
