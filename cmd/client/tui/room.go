package tui

import "github.com/dylanmccormick/ws-chat/internal/protocol"

type Room struct {
	Name             string
	RawMessages      []protocol.Message
	RenderedMessages []string
}

func NewRoom(name string) *Room {
	return &Room{
		Name:             name,
		RawMessages:      []protocol.Message{},
		RenderedMessages: []string{},
	}
}
