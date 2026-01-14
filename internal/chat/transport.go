package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Transport struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
}

// Enviar mensajes formateados
func (t *Transport) Send(header, content string) error {
	msg := fmt.Sprintf("%s|%s\n", header, strings.TrimSpace(content))
	_, err := t.Writer.WriteString(msg)

	if err != nil {
		return err
	}

	return t.Writer.Flush()
}

func (t *Transport) Receive() (string, string, error) {
	line, err := t.Reader.ReadString('\n')

	if err != nil {
		return "", "", err
	}

	parts := strings.SplitN(strings.TrimSpace(line), "|", 2)
	header := parts[0]
	content := ""

	if len(parts) > 1 {
		content = parts[1]
	}

	return header, content, nil
}

func (t *Transport) Close() error {
	return t.Conn.Close()
}
