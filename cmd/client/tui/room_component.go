package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type RoomComponent struct {
	focused bool
	rooms   []string
}

func (cc RoomComponent) View() string {
	str := ""
	for _, room := range cc.rooms {
		str += fmt.Sprintf("%s\n", room)
	}
	return str
}

func (cc RoomComponent) Update(msg tea.Msg) (RoomComponent, tea.Cmd) {
	return cc, nil
}

func (cc RoomComponent) Init() tea.Cmd {
	return nil
}

func NewRoomComponent() *RoomComponent {
	return &RoomComponent{
		focused: true,
		rooms:   []string{"room1", "room2", "lobby"},
	}
}
