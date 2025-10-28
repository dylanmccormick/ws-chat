package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	width  int
	height int

	ChatComponent *ChatComponent
	RoomComponent *RoomComponent
	UserComponent *UserComponent
}

type Component interface {
	Update(msg tea.Msg) (Component, tea.Cmd)
	View() string
	Focus()
	Blur()
}

func Start() {
	rm := NewRootModel()
	p := tea.NewProgram(rm)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func NewRootModel() RootModel {
	return RootModel{
		ChatComponent: NewChatComponent(),
		RoomComponent: NewRoomComponent(),
		UserComponent: NewUserComponent(),
	}
}

func (rm RootModel) Init() tea.Cmd {
	return nil
}

func (rm RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rm.width = msg.Width
		rm.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return rm, tea.Quit
		}
	}
	return rm, nil
}

func (rm RootModel) View() string {
	if rm.height == 0 {
		return "Loading ..."
	}
	rooms := rm.RenderRooms(rm.width/6, rm.height)
	chat := rm.RenderChat(int(float64(rm.width)/float64(1.5)), rm.height)
	users := rm.RenderUsers(rm.width/6, rm.height)
	return lipgloss.JoinHorizontal(lipgloss.Top, rooms, chat, users)
}

func (rm RootModel) RenderChat(width, height int) string {
	content := rm.ChatComponent.View()
	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return headerStyle.Render(content)
}

func (rm RootModel) RenderRooms(width, height int) string {
	content := rm.RoomComponent.View()
	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return headerStyle.Render(content)
}

func (rm RootModel) RenderUsers(width, height int) string {
	content := rm.UserComponent.View()
	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return headerStyle.Render(content)
}
