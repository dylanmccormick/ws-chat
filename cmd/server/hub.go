package server

import (
	"context"
	"log/slog"
)

// The hub is the central event handler
type Hub struct {
	clients map[*User]bool

	messages chan []byte // all inbound messages for the hub. Will have user messages, commands, and announcements
	roomManager *RoomManager
}

// This is the event loop. All messages will come through the hub
func (h *Hub) run() {
	t := Translator{}
	for {
		select {
		case data := <-h.messages:
			message, err := t.BytesToMessage(context.TODO(), data)
			if err != nil {
				slog.Error("Handling json unmarshal error")
				// ignore message
				continue
			}
			h.handleMessage(context.TODO(), message)
		default:
			continue
		}

		// Process user messages

		// Process commands

		// Process announcements
	}
}

func (h *Hub) handleMessage(ctx context.Context, msg Message) {
	switch msg.Typ {
	case "chat":
		room, ok := 


	}
}
