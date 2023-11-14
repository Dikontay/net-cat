package app

import (
	"fmt"
	"net"
	"os"
	"time"
)

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

func writeWelcomeMessage(conn net.Conn) {
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
}

func readName(conn net.Conn) ([]byte, error) {
	nameBuffer := make([]byte, 1024)
	length, err := conn.Read(nameBuffer)
	nameBuffer = nameBuffer[:length-1]
	if err != nil {
		return nil, err
	}
	return nameBuffer, nil
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
