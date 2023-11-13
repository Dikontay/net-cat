package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Name   string
	Writer *bufio.Writer
}

type Message struct {
	Message string
	Sender  *Client
}

var (
	clients   = make(map[*Client]bool)
	serverMux = &sync.Mutex{}
	messages  = make(chan Message)
	history   = []string{}
	defaultPort = "3000"
	connType = "tcp"
	
)



func CheckArgs(args []string) error{
	switch len(args){
	case 2 :
		defaultPort = args[1]
	case 1 :
		defaultPort = "3000"
	default:
		return errors.New("Invalid Usage")
	}
	return  nil
}


func Start(){
	listener, err := net.Listen(connType, ":" + defaultPort)
	defer listener.Close()
	if err != nil {
		fmt.Println("error with listening")
		return
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error with accepting")
			return
		}
		go handleClient(conn)
	}
}


func broadcaster() {
	for {
		msg := <-messages
		for client := range clients {
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

func handleClient(conn net.Conn) {
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
	joinedChat(client)

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
			serverMux.Lock()
			history = append(history, formattedMessage)
			serverMux.Unlock()
			messages <- *msg
		}

	}

	leftChat(client)
	
}

func addClient(client *Client){
		// adding clients
		serverMux.Lock()
		clients[client] = true
		serverMux.Unlock()
}

func deleteClient(client *Client){
	serverMux.Lock()
	delete(clients, client)
	serverMux.Unlock()
}

func joinedChat(client *Client){
	addClient(client)
	messages <- Message{"\n" + client.Name + " has joined the chat.\n", client}
	showHistory(client)
}

func leftChat(client *Client){
	deleteClient(client)
	messages <- Message{"\n" + client.Name + " has left the chat.\n", client}
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

func showHistory(client *Client){
	for _, msg := range history{
	    msg= strings.Trim(msg , "\n")
		client.Writer.WriteString(msg)
		
	}
	client.Writer.Flush()
}
