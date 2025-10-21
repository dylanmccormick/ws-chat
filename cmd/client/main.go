package client

import (
	"fmt"

	"github.com/dylanmccormick/ws-chat/cmd/client/commands"
)

func StartREPL() {
	fmt.Println("starting repl")
	commands.Execute()
}
