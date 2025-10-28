package server

import prot "github.com/dylanmccormick/ws-chat/internal/protocol"

type InternalMessage struct {
	User    *User
	Message prot.Message
}

func NewInternalMessage(user *User, msg prot.Message) *InternalMessage {
	return &InternalMessage{
		User:    user,
		Message: msg,
	}
}

func (m *InternalMessage) EnrichWithUser(user *User) {
	m.User = user
}
