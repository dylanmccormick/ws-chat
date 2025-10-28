package client

import (
	"github.com/dylanmccormick/ws-chat/cmd/client/commands"
	"github.com/dylanmccormick/ws-chat/cmd/client/tui"
)

func StartREPL() {
	commands.Execute()
}

func StartTUI() {
	tui.Start()
}
