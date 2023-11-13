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
	server := app.NewServer(port)
	server.Start()
}
