package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	prot "github.com/dylanmccormick/ws-chat/internal/protocol"
)

func (h *Hub) commandRegisterUser(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
	slog.Info("Registering User", "user", msg.User.username)
	// h.clients[msg.User] = true
	h.clients[msg.User.username] = msg.User
	go reader(msg.User, h.messages)
	rm, err := h.roomManager.GetRoom("lobby")
	if err != nil {
		slog.Error("LOBBY DOES NOT EXIST")
		os.Exit(1)
	}
	rm.Users = append(rm.Users, msg.User)
}

func (h *Hub) commandCreateRoom(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
	slog.Info("User requested to create room", "user", msg.User.username, "room", body.Target)
	// TODO: add a check to make sure a room  doesn't already exist
	h.roomManager.AddRoom(body.Target)
	rm, err := h.roomManager.GetRoom(body.Target)
	if err != nil {
		slog.Error("Was not able to create room", "room", body.Target)
	}
	rm.Users = append(rm.Users, msg.User)
}

func (h *Hub) commandChangeUsername(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
	slog.Info("User requested to change username", "user", msg.User.username, "new_username", body.Target)

	if _, ok := h.clients[body.Target]; ok {
		im := CreateErrorMessage(ctx, "This username is taken")
		out, err := h.translator.MessageToBytes(ctx, im)
		if err != nil {
			slog.Error("Unable to translate message to bytes.", "err", err)
			return
		}
		msg.User.send <- []byte(out)
		return
	}
	slog.Info("Getting user from client map")
	usr := h.clients[msg.User.username]
	slog.Info("retrieved user", "user_is_nil", usr == nil, "username_lookup", msg.User.username)
	slog.Info("Adding new username to client map")
	h.clients[body.Target] = usr
	slog.Info("Updating username in user object")
	usr.username = body.Target
	slog.Info("Deleting user from client map")
	delete(h.clients, msg.User.username)
}

func (h *Hub) commandJoinRoom(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
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
}

func (h *Hub) commandListRoomsForUser(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
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
			Action:   "ListMyRooms",
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
}

func (h *Hub) commandListUsersInRoom(ctx context.Context, msg InternalMessage, body prot.CommandMessage) {
	slog.Info("User requested user information for room", "user", msg.User.username, "room", body.Target)
	r, err := h.roomManager.GetRoom(body.Target)
	if err != nil {
		slog.Error("Error trying to get user information", "room", body.Target, "error", err)
		return
	}
	users := []string{}
	for _, user := range r.Users {
		users = append(users, user.username)
	}
	data, err := json.Marshal(users)
	if err != nil {
		slog.Error("Unable to create response data", "err", err)
	}
	sendMsg := prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Target:   r.Name,
			Type:     "commandResponse",
			Action:   "ListRoomUsers",
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
}
