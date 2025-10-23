package main

import (
	"fmt"
	"os"

	// "github.com/dylanmccormick/ws-chat/cmd/server"
	// "github.com/dylanmccormick/ws-chat/cmd/client"
	"github.com/dylanmccormick/ws-chat/cmd/ws-chat"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ws-chat [start|repl|tui]")
	}

	wschat.Execute()
	// switch os.Args[1]{
	// case "start":
	// 	server.StartServer()
	// case "repl":
	// 	client.StartREPL()
	// default:
	// 	fmt.Printf("Unknown Command %s\n", os.Args[1])
	// }
}

