package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	prot "github.com/dylanmccormick/ws-chat/internal/protocol"
	"github.com/gorilla/websocket"
)

// The hub is the central event handler
type Hub struct {
	// clients    map[*User]bool
	clients    map[string]*User
	register   chan *User
	unregister chan *User

	messages    chan InternalMessage // all inbound messages for the hub. Will have user messages, commands, and announcements
	roomManager *RoomManager
	translator  Translator
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*User),
		messages:    make(chan InternalMessage, 10),
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

func (h *Hub) handleMessage(ctx context.Context, intMsg InternalMessage) {
	msg := intMsg.Message
	slog.Info("Got message with body type", "type", reflect.TypeOf(msg.Body))
	switch body := msg.Body.(type) {
	case prot.ChatMessage:
		h.handleChat(ctx, intMsg, body)
	case prot.AnnouncementMessage:
		h.handleAnnouncement(ctx, intMsg, body)
	case prot.ErrorMessage:
		h.handleError(ctx, intMsg, body)
	case prot.CommandMessage:
		slog.Info("Handling command")
		h.handleCommand(ctx, intMsg, body)
	}
}

func (h *Hub) handleChat(ctx context.Context, msg InternalMessage, body prot.ChatMessage) {
	// TODO: Assert body is of chatMessage type
	room, err := h.roomManager.GetRoom(body.Target)
	if err != nil {
		slog.Error("Unable to resolve target for chat message", "message", msg, "body", body)
		return
	}
	body.UserName = msg.User.username
	msg.Message.Body = body
	data, err := h.translator.MessageToBytes(ctx, msg)
	if err != nil {
		slog.Error("Unable to convert message to bytes", "message", msg)
		return
	}
	h.broadcast(ctx, data, room)
}

func (h *Hub) handleAnnouncement(ctx context.Context, msg InternalMessage, body prot.AnnouncementMessage) {
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

func (h *Hub) handleError(ctx context.Context, msg InternalMessage, body prot.ErrorMessage) {
	slog.Warn("Error not yet implemented. This is all you get buddy", "message", body.Message)
}

func (h *Hub) handleCommand(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
	switch body.Action {
	case "RegisterUser":
		h.commandRegisterUser(ctx, msg, body)
	case "CreateRoom":
		h.commandCreateRoom(ctx, msg, body)
	case "JoinRoom":
		h.commandJoinRoom(ctx, msg, body)
	case "ListMyRooms":
		h.commandListRoomsForUser(ctx, msg, body)
	case "ListRoomUsers":
		h.commandListUsersInRoom(ctx, msg, body)
		slog.Info("User requested user information", "user", msg.User.username, "room", body.Target)

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

	msg := prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Action:   "RegisterUser",
			UserName: u.username,
		},
	}
	internalMessage := InternalMessage{
		Message: msg,
		User:    u,
	}
	slog.Info("Posting test message to messages queue", "message", msg)
	h.messages <- internalMessage
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

func userInRoom(room *Room, targetUser *User) bool {
	return slices.Contains(room.Users, targetUser)
}
