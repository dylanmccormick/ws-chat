package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/gorilla/websocket"
)

// The hub is the central event handler
type Hub struct {
	clients    map[*User]bool
	register   chan *User
	unregister chan *User

	messages    chan Message // all inbound messages for the hub. Will have user messages, commands, and announcements
	roomManager *RoomManager
	translator  Translator
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*User]bool),
		messages:    make(chan Message),
		register:    make(chan *User),
		unregister:  make(chan *User),
		roomManager: NewRoomManager(),
		translator:  Translator{},
	}
}

// This is the event loop. All messages will come through the hub
func (h *Hub) run() {
	slog.Info("Starting hub")
	h.roomManager.AddRoom("lobby")
	for {
		select {
		case client := <-h.register:
			go h.registerClient(client)
		case message := <-h.messages:
			slog.Info("Received a message", "msg", message)
			h.handleMessage(context.TODO(), message)
		default:
			continue
		}
	}
}

func (h *Hub) handleMessage(ctx context.Context, msg Message) {
	slog.Info("Got message with body type", "type", reflect.TypeOf(msg.Body))
	switch body := msg.Body.(type) {
	case ChatMessage:
		h.handleChat(ctx, msg, body)
	case AnnouncementMessage:
		h.handleAnnouncement(ctx, msg, body)
	case ErrorMessage:
		h.handleError(ctx, msg, body)
	case CommandMessage:
		slog.Info("Handling command")
		h.handleCommand(ctx, msg, body)
	}
}

func (h *Hub) handleChat(ctx context.Context, msg Message, body ChatMessage) {
	// TODO: Assert body is of chatMessage type
	room, err := h.roomManager.GetRoom(body.Target)
	if err != nil {
		slog.Error("Unable to resolve target for chat message", "message", msg, "body", body)
		return
	}
	data, err := h.translator.MessageToBytes(ctx, msg)
	if err != nil {
		slog.Error("Unable to convert message to bytes", "message", msg)
		return
	}
	h.broadcast(ctx, data, room)
}

func (h *Hub) handleAnnouncement(ctx context.Context, msg Message, body AnnouncementMessage) {
	room, err := h.roomManager.GetRoom(body.Target)
	if err != nil {
		slog.Error("Unable to resolve target for announcement message", "message", msg, "body", body)
		return
	}
	data, err := h.translator.MessageToBytes(ctx, msg)
	if err != nil {
		slog.Error("Unable to convert message to bytes", "message", msg)
		return
	}
	h.broadcast(ctx, data, room)
}

func (h *Hub) handleError(ctx context.Context, msg Message, body ErrorMessage) {}

func (h *Hub) handleCommand(ctx context.Context, msg Message, body CommandMessage) {
	switch body.Action {
	case "RegisterUser":
		slog.Info("Registering User", "user", msg.User.username)
		h.clients[msg.User] = true
		go reader(msg.User, h.messages)
		rm, err := h.roomManager.GetRoom("lobby")
		if err != nil {
			slog.Error("LOBBY DOES NOT EXIST")
			os.Exit(1)
		}
		rm.Users = append(rm.Users, msg.User)
	case "CreateRoom":
		slog.Info("User requested to create room", "user", msg.User.username, "room", body.Target)
		h.roomManager.AddRoom(body.Target)
		rm, err := h.roomManager.GetRoom(body.Target)
		if err != nil {
			slog.Error("Was not able to create room", "room", body.Target)
		}
		rm.Users = append(rm.Users, msg.User)
	case "JoinRoom":
		slog.Info("User requested to join room", "user", msg.User.username, "room", body.Target)
		rm, err := h.roomManager.GetRoom(body.Target)
		if err != nil {
			slog.Error("Was not able to join room", "room", body.Target)
		}
		rm.Users = append(rm.Users, msg.User)
		msg := Message{
			Typ:  "announcement",
			User: msg.User,
			Body: AnnouncementMessage{
				Message: fmt.Sprintf("User %s has joined the room", msg.User.username),
				Target:  body.Target,
			},
		}
		h.messages <- msg
	default:
		slog.Warn("Received command with unexpected action", "action", body.Action)
	}
}

func (h *Hub) broadcast(ctx context.Context, data []byte, room *Room) {
	// TODO: Follow gorilla pattern of dropping bad users
	for _, u := range room.Users {
		u.send <- data
	}
}

func (h *Hub) registerClient(u *User) {
	h.promptForUsername(u)

	msg := Message{
		Typ:  "command",
		User: u,
		Body: CommandMessage{
			Action: "RegisterUser",
		},
	}
	slog.Info("Posting test message to messages queue", "message", msg)
	h.messages <- msg
}

func (h *Hub) promptForUsername(u *User) {
	var validUsername bool
	validUsername = false
	var username string
	conn := u.conn
	for !validUsername {
		WriteToConn(conn, []byte("Please send your username"))

		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("That username is bogus", "error", err)
		}
		slog.Info("Got username", "username", message)

		username = string(bytes.TrimSpace(bytes.ReplaceAll(message, []byte("\n"), []byte(" "))))
		validUsername = true
	}
	WriteToConn(conn, fmt.Appendf(nil, "Welcome to the lobby, %s", username))
	u.username = username
}

func WriteToConn(conn *websocket.Conn, message []byte) {
	ws, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error("An error occurred with NextWriter: ", "error", err)
	}
	ws.Write(message)
	ws.Close()
}
