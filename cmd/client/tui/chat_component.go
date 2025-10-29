package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SendChatMessage struct {
	Message string
}

type ChatComponent struct {
	focused bool
	input   textinput.Model
}

func (cc *ChatComponent) ViewRoom(room *Room) string {
	return cc.ViewMessages(room.RenderedMessages)
}

func (cc *ChatComponent) ViewMessages(messages []string) string {
	str := ""
	for _, msg := range messages {
		str += msg + "\n"
	}
	return str
}

func (cc ChatComponent) View() string {
	str := ""
	return str
}

func (cc *ChatComponent) Update(msg tea.Msg, room string) (ChatComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			text := cc.input.Value()
			cc.input.SetValue("")
			return *cc, func() tea.Msg {
				return SendChatMessage{Message: text}
			}
		}
	}
	var cmd tea.Cmd
	cc.input, cmd = cc.input.Update(msg)
	return *cc, cmd
}

func (cc *ChatComponent) Init() tea.Cmd {
	return nil
}

func NewChatComponent() *ChatComponent {
	ti := textinput.New()
	ti.Width = 100
	ti.Focus()
	return &ChatComponent{
		focused: true,
		input:   ti,
	}
}

func (cc *ChatComponent) Focus() {
	cc.focused = true
}

func (cc *ChatComponent) Blur() {
	cc.focused = false
}
