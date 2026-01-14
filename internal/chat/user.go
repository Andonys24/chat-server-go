package chat

import (
	"net"
	"time"
)

type User struct {
	Username  string
	LoginTime time.Time
	Transport *Transport // Manejador de protocolo
	Conn      net.Conn
}

func NewUser(username string, conn net.Conn, transport *Transport) *User {
	return &User{
		Username:  username,
		LoginTime: time.Now(),
		Transport: transport,
		Conn:      conn,
	}
}
