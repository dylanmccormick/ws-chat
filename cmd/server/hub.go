package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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
		slog.Info("Registering User", "user", msg.User.username)
		// h.clients[msg.User] = true
		go reader(msg.User, h.messages)
		rm, err := h.roomManager.GetRoom("lobby")
		if err != nil {
			slog.Error("LOBBY DOES NOT EXIST")
			os.Exit(1)
		}
		rm.Users = append(rm.Users, msg.User)
	case "CreateRoom":
		slog.Info("User requested to create room", "user", msg.User.username, "room", body.Target)
		// TODO: add a check to make sure a room  doesn't already exist
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
		// TODO: add a check to make sure the user is not already in the room
		if userInRoom(rm, msg.User) {
			slog.Warn("User already in room", "room", rm.Name, "user", msg.User.username)
			im := CreateErrorMessage(ctx, "This user already is in this room")
			h.messages <- im
			return
		}

		rm.Users = append(rm.Users, msg.User)
		sendMsg := prot.Message{
			Typ: "announcement",
			Body: prot.AnnouncementMessage{
				Message:  fmt.Sprintf("User %s has joined the room", msg.User.username),
				Target:   body.Target,
				UserName: msg.User.username,
			},
		}
		intMsg := InternalMessage{
			User:    msg.User,
			Message: sendMsg,
		}
		h.messages <- intMsg
	case "ListMyRooms":
		slog.Info("User requested room information", "user", msg.User.username)
		rooms := h.roomManager.ListRooms()
		userRooms := []string{}
		for _, room := range rooms {
			r, err := h.roomManager.GetRoom(room)
			if err != nil {
				slog.Error("Error trying to get room information", "room", body.Target, "error", err)
			}
			for _, user := range r.Users {
				if user.username == msg.User.username {
					userRooms = append(userRooms, r.Name)
				}
			}
		}
		data, err := json.Marshal(userRooms)
		if err != nil {
			slog.Error("Unable to create response data", "err", err)
			return
		}
		sendMsg := prot.Message{
			Typ: "command",
			Body: prot.CommandMessage{
				Target:   msg.User.username,
				Type:     "commandResponse",
				Action:   "",
				Data:     data,
				UserName: msg.User.username,
			},
		}
		intMsg := InternalMessage{
			User:    msg.User,
			Message: sendMsg,
		}
		out, err := h.translator.MessageToBytes(ctx, intMsg)
		if err != nil {
			slog.Error("Unable to translate message to bytes.", "err", err)
			return
		}
		msg.User.send <- []byte(out)

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
