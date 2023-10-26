package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn  net.Conn
	name  string
	msgCh chan string
}

var (
	HOST      = "localhost"
	PORT      = "3000"
	TYPE      = "tcp"
	clients   = make(map[*Client]bool)
	clientMux sync.Mutex
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	handleErrors(err)
	defer listen.Close()
	for {
		conn, err := listen.Accept()

		handleErrors(err)
		client := &Client{
			conn:  conn,
			msgCh: make(chan string),
		}
		go handleIncomingRequest(client)

		// conn.Close()
	}
}

func handleIncomingRequest(client *Client) {
	defer client.conn.Close()
	client.conn.Write([]byte("Please enter your name: "))
	name := []byte{}
	for {
		buffer := make([]byte, 1)
		n, err := client.conn.Read(buffer)
		if err != nil {
			handleErrors(err)
			break
		}
		if n == 0 {
			// No data was read; continue to read.
			continue
		}

		if buffer[0] == '\n' {
			// If a newline character is encountered, stop reading.
			break
		}

		name = append(name, buffer[0])
	}

	client.name = string(name)
	clientMux.Lock()
	clients[client] = true
	clientMux.Unlock()
	// client.conn.Write([]byte(fmt.Sprintf("%s has joined the caht", client.name)))
	broadcast(fmt.Sprintf("%s has joined the chat", client.name), &Client{})
	go client.receiveMessages()

	for msg := range client.msgCh {
		_, err := client.conn.Write([]byte(msg))
		if err != nil {
			break
		}
	}
	clientMux.Lock()
	delete(clients, client)
	clientMux.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat", client.name), &Client{})
}

func (client *Client) receiveMessages() {
	// scanner := bufio.NewScanner(client.conn)



	for {
		timestap := time.Now().Format("02-Jan-06 15:04:05 MST")
		// fmt.Println(buffer)
		client.conn.Write([]byte(fmt.Sprintf("[%s][%s]: ", string(client.name), timestap)))
		buffer := make([]byte, 1024)
		_, err := client.conn.Read(buffer)
		if err != nil {
			broadcast(fmt.Sprintf("%s has left the chat", client.name), &Client{})
			return
		}
		
		broadcast(fmt.Sprintf("[%s][%s]:%s\n", client.name, timestap,  string(buffer)), client)
	}

	
}

func broadcast(message string, currentClient *Client) {
	clientMux.Lock()
	defer clientMux.Unlock()
	for client := range clients {
		if client != currentClient {
			select {
			case client.msgCh <- message:
				// fmt.Println("send messages")
			default:
				// fmt.Printf("Failed to send message to client %s\n", client.name)
			}
		}
	}
	//close(currentClient.msgCh)
}

func handleErrors(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
