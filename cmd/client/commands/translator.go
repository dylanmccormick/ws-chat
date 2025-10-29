package commands

import (
	"encoding/json"
	"log/slog"

	"github.com/dylanmccormick/ws-chat/internal/protocol"
)

type Translator struct{}

func (t *Translator) BytesToMessage(data []byte) (protocol.Message, error) {
	var msg protocol.Message
	err := msg.UnmarshalJSON(data)
	if err != nil {
		slog.Error("failed to unmarshal json from message", "message", string(data), "user", "user") // TODO: Add user to context
		return CreateErrorMessage("Unable to parse message from user"), nil
	}
	return msg, nil
}

func (t *Translator) MessageToBytes(protocol.Message) ([]byte, error) {
	msg := protocol.Message{}
	data, err := msg.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal json from Message", "message", msg)
	}
	return data, nil
}

func CreateErrorMessage(msg string) protocol.Message {
	errMsg, err := json.Marshal(protocol.ErrorMessage{Message: msg})
	if err != nil {
		// if this doesn't work that means I have crafted a bad message... which means the program is bad
		panic(err)
	}
	return protocol.Message{
		Typ:  "error",
		Body: errMsg,
	}
}
