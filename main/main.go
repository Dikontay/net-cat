package main

import (
	"fmt"
	"netcat/cmd/app"
	"os"
)

var usage = "[USAGE]: ./TCPChat $port"

func main() {
	err := app.CheckArgs(os.Args)
	if err != nil {
		fmt.Println(usage)
		return
	}
	app.Start()
}
