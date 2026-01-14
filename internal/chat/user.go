package chat

import (
	"net"
)

type User struct {
	Username  string
	Transport *Transport // Manejador de protocolo
	Conn      net.Conn
}

func NewUser(username string, conn net.Conn, transport *Transport) *User {
	return &User{
		Username:  username,
		Transport: transport,
		Conn:      conn,
	}
}
