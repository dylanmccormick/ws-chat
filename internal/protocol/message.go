package protocol

import "encoding/json"

// This is the shared message col between the different layers of the application.
// Both the client and the server rely on this col to create and decode messages.
type Message struct {
	Typ  string `json:"type"`
	Body any    `json:"-"`
}

type ChatMessage struct {
	Message  string `json:"message"`
	Target   string `json:"target"` // The room for the chat message
	UserName string `json:"username,omitempty"`
}

type AnnouncementMessage struct {
	Message  string `json:"message"`
	Target   string `json:"target"` // The room for the announcement message
	UserName string `json:"username,omitempty"`
}

type ErrorMessage struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	UserName string `json:"username,omitempty"`
}

type CommandMessage struct {
	Target   string          `json:"target"`
	Type     string          `json:"command"`
	Action   string          `json:"action"`
	UserName string          `json:"username,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
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
		var commandBody CommandMessage
		if err := json.Unmarshal(temp.Body, &commandBody); err != nil {
			return err
		}
		m.Body = commandBody
	case "error":
		var errorBody ErrorMessage
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
