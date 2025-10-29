package server

import (
	"context"
	"encoding/json"
	"log/slog"

	prot "github.com/dylanmccormick/ws-chat/internal/protocol"
)

// This is the translator. Also known as the byte wizard. It translates []byte into Message{} and takes in some user context to be able to ascribe a Message to a User

type Translator struct{}

func (t *Translator) BytesToMessage(ctx context.Context, data []byte) (InternalMessage, error) {
	var msg prot.Message
	err := msg.UnmarshalJSON(data)
	if err != nil {
		slog.Error("failed to unmarshal json from message", "message", string(data), "user", "user") // TODO: Add user to context
		return CreateErrorMessage(ctx, "Unable to parse message from user"), nil
	}
	return InternalMessage{
		Message: msg,
	}, nil
}

func (t *Translator) MessageToBytes(ctx context.Context, internalMsg InternalMessage) ([]byte, error) {
	msg := internalMsg.Message
	// TODO: This might need an internal message Marshal Function
	data, err := msg.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal json from Message", "message", msg)
	}
	return data, nil
}

func CreateErrorMessage(ctx context.Context, msg string) InternalMessage {
	errMsg, err := json.Marshal(prot.ErrorMessage{Message: msg})
	if err != nil {
		// if this doesn't work that means I have crafted a bad message... which means the program is bad
		panic(err)
	}
	return InternalMessage{
		User: &User{}, // TODO: Get user from context
		Message: prot.Message{
			Typ:  "error",
			Body: errMsg,
		},
	}
}

