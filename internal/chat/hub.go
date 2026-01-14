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
}

// Estructura auxiliar para pedir registro
type RegisterRequest struct {
	Nickname  string
	Conn      net.Conn
	Transport *Transport
}

func NewHub(maxConn int) *Hub {
	return &Hub{
		Clients:        make(map[string]*User),
		MaxConnections: maxConn,
		Register:       make(chan *RegisterRequest),
		Unregister:     make(chan net.Conn),
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub de usuarios iniciado...")
	for {
		select {
		case <-h.Register:
		case <-h.Unregister:
		}
	}
}
