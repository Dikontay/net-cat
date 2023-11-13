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

func CheckArgs(args []string) string {
	var port string
	switch len(args) {
	case 2:
		port = args[1]
	case 1:
		port = "3000"
	default:
		return ""
	}
	return port
}

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
	var welcome string = `
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`
	_, err := conn.Write([]byte(welcome))
	if err != nil {
		fmt.Println("Error writing message")
		os.Exit(1)
	}
	_, err = conn.Write([]byte("please enter your name : "))
	if err != nil {
		fmt.Println("Error writing message")
		os.Exit(1)
	}
	nameBuffer := make([]byte, 1024)
	length, err := conn.Read(nameBuffer)
	nameBuffer = nameBuffer[:length-1]
	if err != nil {
		fmt.Println("Error reading user name")
		os.Exit(1)
	}
	writer := bufio.NewWriter(conn)
	client := &Client{
		Name:   string(nameBuffer),
		Writer: writer,
	}
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
			mux.Lock()
			s.history = append(s.history, formattedMessage)
			mux.Unlock()
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

func (s *Server) joinedChat(client *Client) {
	s.addClient(client)
	s.messages <- Message{"\n" + client.Name + " has joined the chat.\n", client}
	s.showHistory(client)
}

func (s *Server) leftChat(client *Client) {
	s.deleteClient(client)
	s.messages <- Message{"\n" + client.Name + " has left the chat.\n", client}
}

func sendPrompt(client *Client) {
	timestamp := time.Now().Format("02-Jan-06 15:04:05 MST")
	_, err := client.Writer.WriteString(fmt.Sprintf("[%s][%s]:", timestamp, client.Name))
	if err != nil {
		fmt.Println("Error writing string")
		os.Exit(1)
	}
	err = client.Writer.Flush()
	if err != nil {
		fmt.Println("Error flushing")
		os.Exit(1)
	}
}

func (s *Server) showHistory(client *Client) {
	for _, msg := range s.history {
		msg = strings.Trim(msg, "\n")
		client.Writer.WriteString(msg)

	}
	client.Writer.Flush()
}
