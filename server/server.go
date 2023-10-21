package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var (
	HOST = "localhost"
	PORT = "3000"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	handleErrors(err)
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		handleErrors(err)
		go handleIncomingRequest(conn)

	}

}

func handleIncomingRequest(conn net.Conn) {
	//store incoming data
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	handleErrors(err)
	//respond
	time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	//fmt.Println(buffer)
	conn.Write([]byte("Hi back!"))
	conn.Write([]byte(time))
	conn.Write(buffer)
	conn.Close()
}

func handleErrors(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
