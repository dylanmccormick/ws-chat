package tui

import (
	"bytes"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylanmccormick/ws-chat/cmd/client/commands"
	"github.com/dylanmccormick/ws-chat/internal/protocol"
	"github.com/gorilla/websocket"
)

type RootModel struct {
	width  int
	height int

	roomsMap    map[string]*Room
	CurrentRoom *Room

	ChatComponent *ChatComponent
	RoomComponent *RoomComponent
	UserComponent *UserComponent

	Conn *websocket.Conn
	sub  chan protocol.Message

	MessageCount int
	ChatsSent    int
}

type Component interface {
	Update(msg tea.Msg) (Component, tea.Cmd)
	View() string
	Focus()
	Blur()
}

func Start() {
	conn := commands.CreateConnection()
	rm := NewRootModel(conn)
	p := tea.NewProgram(rm)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func NewRootModel(conn *websocket.Conn) RootModel {
	lobby := NewRoom("lobby")
	return RootModel{
		CurrentRoom:   lobby,
		roomsMap:      map[string]*Room{"lobby": lobby},
		ChatComponent: NewChatComponent(),
		RoomComponent: NewRoomComponent(),
		UserComponent: NewUserComponent(),
		Conn:          conn,
		sub:           make(chan protocol.Message, 10),
		MessageCount:  0,
		ChatsSent:     0,
	}
}

func (rm RootModel) Init() tea.Cmd {
	return tea.Batch(
		ListenForMessages(rm.Conn, rm.sub),
		ReceiveMessage(rm.sub),
	)
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
	case protocol.Message:
		rm.ProcessMessage(msg)
		return rm, ReceiveMessage(rm.sub)
	case SendChatMessage:
		rm.ChatsSent++
		return rm, rm.SendChatMessage(msg)
	}

	var cmd tea.Cmd
	*rm.ChatComponent, cmd = rm.ChatComponent.Update(msg, rm.CurrentRoom.Name)

	return rm, cmd
}

func (rm RootModel) View() string {
	if rm.height == 0 {
		return "Loading ..."
	}
	rooms := rm.RenderRooms(rm.width/6, rm.height-4)
	chat := rm.RenderChat(int(float64(rm.width)/float64(1.5)), rm.height-4)
	users := rm.RenderUsers(rm.width/6, rm.height-4)
	page := lipgloss.JoinHorizontal(lipgloss.Top, rooms, chat, users)
	footer := rm.RenderFooter(rm.width, 4)
	return lipgloss.JoinVertical(lipgloss.Left, page, footer)
}

func (rm RootModel) RenderFooter(width, height int) string {
	content := fmt.Sprintf("Current room: %s, messages: %d, chats_sent: %d", rm.CurrentRoom.Name, rm.MessageCount, rm.ChatsSent)
	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
	return headerStyle.Render(content)
}

func (rm RootModel) RenderChat(width, height int) string {
	messages := rm.ChatComponent.ViewRoom(rm.CurrentRoom)
	input := rm.ChatComponent.input.View()

	chatStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	inputStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(4).
		Foreground(lipgloss.Color("64"))

	return lipgloss.JoinVertical(lipgloss.Left, chatStyle.Render(messages), inputStyle.Render(input))
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

func ListenForMessages(c *websocket.Conn, sub chan protocol.Message) tea.Cmd {
	translator := commands.Translator{}
	return func() tea.Msg {
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}
			data = bytes.TrimSpace(bytes.ReplaceAll(data, []byte("\n"), []byte(" ")))
			msg, err := translator.BytesToMessage(data)
			if err != nil {
				panic(err)
			}

			sub <- msg
		}
	}
}

func (rm *RootModel) SendChatMessage(msg SendChatMessage) tea.Cmd {
	return func() tea.Msg {
		bytMsg := commands.CreateChatMessage(msg.Message, rm.CurrentRoom.Name)
		err := rm.Conn.WriteMessage(websocket.TextMessage, bytMsg)
		if err != nil {
			fmt.Printf("We got an error writing: %s", err)
			panic(err)
		}
		return nil
	}
}

func ReceiveMessage(sub chan protocol.Message) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (rm RootModel) ProcessMessage(msg protocol.Message) {
	rm.MessageCount++
	switch body := msg.Body.(type) {
	case protocol.ChatMessage:
		roomName := body.Target
		room, ok := rm.roomsMap[roomName]
		if !ok {
			// This might need to be an error but I'm not sure how to show those yet
			return
		}
		room.RawMessages = append(room.RawMessages, msg)
		room.RenderedMessages = append(room.RenderedMessages, renderChat(body))
		return

	case protocol.AnnouncementMessage:
		roomName := body.Target
		room, ok := rm.roomsMap[roomName]
		if !ok {
			// This might need to be an error but I'm not sure how to show those yet
			return
		}
		room.RawMessages = append(room.RawMessages, msg)
		room.RenderedMessages = append(room.RenderedMessages, renderAnnouncement(body))
		return
	}
}

func renderChat(msg protocol.ChatMessage) string {
	return fmt.Sprintf("%s: %s", msg.UserName, msg.Message)
}

func renderAnnouncement(msg protocol.AnnouncementMessage) string {
	return fmt.Sprintf("%s", msg.Message)
}
