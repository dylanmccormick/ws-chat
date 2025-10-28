package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	prot "github.com/dylanmccormick/ws-chat/internal/protocol"
)

type ChatComponent struct {
	focused  bool
	messages []prot.ChatMessage

	input textinput.Model
}

func (cc ChatComponent) View() string {
	str := ""
	for _, msg := range cc.messages {
		str += fmt.Sprintf("%s: %s\n", msg.UserName, msg.Message)
	}
	return str
}

func (cc ChatComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return cc, nil
}

func (cc ChatComponent) Init() tea.Cmd {
	return nil
}

func NewChatComponent() *ChatComponent {
	return &ChatComponent{
		focused: true,
		messages: []prot.ChatMessage{
			{
				Message:  "this is test message 1",
				UserName: "bozo",
			},
			{
				Message:  "this is test message 2",
				UserName: "bingo",
			},
			{
				Message:  "this is test message 3",
				UserName: "alfonz",
			},
		},
	}
}
