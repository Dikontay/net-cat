package app

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var mux = &sync.Mutex{}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":"+s.Port)
	defer listener.Close()
	if err != nil {
		fmt.Println("error with listening")
		return
	}
	go s.broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error with accepting")
			return
		}
		go s.handleClient(conn)
	}
}

func (s *Server) broadcaster() {
	for {
		msg := <-s.messages
		for client := range s.clients {
			if msg.Sender != client {
				_, err := client.Writer.WriteString(msg.Message)
				if err != nil {
					fmt.Println("Error broadcasting")
					os.Exit(1)
				}
				err = client.Writer.Flush()
				if err != nil {
					fmt.Println("Error flushing")
					os.Exit(1)
				}
				sendPrompt(client)
			}
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	if len(s.clients) == 10 {
		conn.Write([]byte("Maximum 10 connections available\n"))
		return
	}

	writeWelcomeMessage(conn)
	nameBuffer, err := readName(conn)
	if err != nil {
		fmt.Println("error reading a name")
		return
	}
	client := newClient(nameBuffer, conn)
	s.joinedChat(client)

	reader := bufio.NewReader(conn)
	for {
		sendPrompt(client)
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected.\n", client.Name)
			} else {
				fmt.Printf("Error reading from client %s: %s\n", client.Name, err)
			}
			break
		}
		message = strings.Trim(message, "\r\n") // Trim the message

		if message != "" {

			formattedMessage := fmt.Sprintf("\n[%s] [%s]: %s\n", time.Now().Format("02-Jan-06 15:04:05 MST"), nameBuffer, message)
			msg := &Message{Message: formattedMessage, Sender: client}
			s.addToHistory(formattedMessage)
			s.messages <- *msg
		}

	}

	s.leftChat(client)
}

func (s *Server) addClient(client *Client) {
	// adding clients
	mux.Lock()
	s.clients[client] = true
	mux.Unlock()
}

func (s *Server) deleteClient(client *Client) {
	mux.Lock()
	delete(s.clients, client)
	mux.Unlock()
}

func (s *Server) addToHistory(message string) {
	mux.Lock()
	s.history = append(s.history, message)
	mux.Unlock()
}

func (s *Server) joinedChat(client *Client) {
	s.addClient(client)
	s.messages <- Message{"\n" + client.Name + " has joined the chat.\n", client}
	s.showHistory(client)
}

func (s *Server) leftChat(client *Client) {
	s.deleteClient(client)
	s.messages <- Message{"\n" + client.Name + " has left the chat.\n", client}
}

func (s *Server) showHistory(client *Client) {
	for _, msg := range s.history {
		msg = strings.Trim(msg, "\n")
		_, err := client.Writer.WriteString(msg + "\n")
		if err != nil {
			fmt.Println("error writing string")
			return
		}
		err = client.Writer.Flush()
		if err != nil {
			fmt.Println("error flushing")
			return
		}

	}

}
