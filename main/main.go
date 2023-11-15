package main

import (
	"fmt"
	"netcat/cmd/app"
	"os"
)

const usage = "[USAGE]: ./TCPChat $port"

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
