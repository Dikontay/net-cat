package main

import (
	"fmt"
	"net"
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
		//conn.Close()
	}

}

func handleIncomingRequest(conn net.Conn) {
	//store incoming data
	name := make([]byte, 1024)
	conn.Write([]byte("Please enter your name:"))
	_, err := conn.Read(name)
	handleErrors(err)
	//fmt.Print(string(name))
	//fmt.Print(name)
	for {
		timestap := time.Now().Format("02-Jan-06 15:04:05 MST")
		//fmt.Println(buffer)
		conn.Write([]byte(fmt.Sprintf("%s %s: ", string(name), timestap)))
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Printf("User %s has left.\n", string(name))
			break
		}

	}
	conn.Close()
	//respond

}

func handleErrors(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
