package app

import (
	"bufio"
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
