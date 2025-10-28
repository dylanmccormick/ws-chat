package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type UserComponent struct {
	focused bool
	users   []string
}

func (cc UserComponent) View() string {
	str := ""
	for _, user := range cc.users {
		str += fmt.Sprintf("%s\n", user)
	}
	return str
}

func (cc UserComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return cc, nil
}

func (cc UserComponent) Init() tea.Cmd {
	return nil
}

func NewUserComponent() *UserComponent {
	return &UserComponent{
		focused: true,
		users:   []string{"bozo", "bingo", "alfonz"},
	}
}
