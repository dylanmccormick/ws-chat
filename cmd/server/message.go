package server

type Message struct {
	Typ      string         `json:"type"`
	User     *User          `json:"-"` // We will pull from the context. That way user can't imitate another user
	Body     any            `json:"-"`
	Metadata map[string]any `json:"-"` // This field is just for us. The client need not know about it
}

type ChatMessage struct {
	Message string `json:"message"`
	Target  string `json:"target"` // The room for the chat message
}

type AnnouncementMessage struct {
	Message string `json:"message"`
	Target  string `json:"target"` // The room for the announcement message
}

type ErrorMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type CommandMessage struct {
	Target string `json:"target"`
	Type   string `json:"command"`
	Action string `json:"action"`
}

func (m *Message) EnrichWithUser(user *User) {
	m.User = user
}
