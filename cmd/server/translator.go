package server

import (
	"context"
	"encoding/json"
	"log/slog"
)

// This is the translator. Also known as the byte wizard. It translates []byte into Message{} and takes in some user context to be able to ascribe a Message to a User

type Translator struct{}

func (t *Translator) BytesToMessage(ctx context.Context, data []byte) (Message, error) {
	var msg Message
	err := msg.UnmarshalJSON(data)
	if err != nil {
		slog.Error("failed to unmarshal json from message", "message", string(data), "user", "user") // TODO: Add user to context
		return CreateErrorMessage(ctx, "Unable to parse message from user"), nil
	}
	return msg, nil
}

func (t *Translator) MessageToBytes(ctx context.Context, msg Message) ([]byte, error) {
	data, err := msg.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal json from Message", "message", msg)
	}
	return data, nil
}

func CreateErrorMessage(ctx context.Context, msg string) Message {
	errMsg, err := json.Marshal(errorMessage{Message: msg})
	if err != nil {
		// if this doesn't work that means I have crafted a bad message... which means the program is bad
		panic(err)
	}
	return Message{
		Typ:  "error",
		User: &User{}, // TODO: Get user from context
		Body: errMsg,
	}
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Typ = temp.Type

	switch temp.Type {
	case "chat":
		var chatBody ChatMessage
		if err := json.Unmarshal(temp.Body, &chatBody); err != nil {
			return err
		}
		m.Body = chatBody
	case "command":
		var commandBody commandMessage
		if err := json.Unmarshal(temp.Body, &commandBody); err != nil {
			return err
		}
		m.Body = commandBody
	case "error":
		var errorBody errorMessage
		if err := json.Unmarshal(temp.Body, &errorBody); err != nil {
			return err
		}
		m.Body = errorBody
	}

	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	var temp struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}
	temp.Type = m.Typ

	body, err := json.Marshal(m.Body)
	if err != nil {
		return []byte{}, err
	}
	temp.Body = body

	return json.Marshal(temp)
}
