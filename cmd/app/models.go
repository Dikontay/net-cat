package app

import (
	"bufio"
	"net"
)

type Server struct {
	Port string

	clients  map[*Client]bool
	messages chan Message
	history  []string
}

type Client struct {
	Name   string
	Writer *bufio.Writer
}

type Message struct {
	Message string
	Sender  *Client
}

// NewServer creates a new Server instance with default values
func NewServer(port string) *Server {
	return &Server{
		Port:     port,
		clients:  make(map[*Client]bool),
		messages: make(chan Message),
		history:  []string{},
	}
}

func newClient(nameBuffer []byte, conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	client := &Client{
		Name:   string(nameBuffer),
		Writer: writer,
	}
	return client
}
