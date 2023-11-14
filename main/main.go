package main

import (
	"fmt"
	"netcat/cmd/app"
	"os"
)

var usage = "[USAGE]: ./TCPChat $port"

func main() {
	port := app.CheckArgs(os.Args)
	if port == "" {
		fmt.Println(usage)
		return
	}
	fmt.Printf("Starting server on port %s\n", port)
	server := app.NewServer(port)
	server.Start()
}
