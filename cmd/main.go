package main

import (
	"fmt"
	"os"

	"github.com/dylanmccormick/ws-chat/cmd/ws-chat"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ws-chat [start|repl|tui]")
	}

	wschat.Execute()
}

