package main

import (
	"fmt"
	"github.com/dylanmccormick/ws-chat/cmd/server"
)


func main() {
	fmt.Println("Hello and welcome to websocket chat")

	server.StartServer()
}

