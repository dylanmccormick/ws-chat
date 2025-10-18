package server

type Message struct {
	Typ  string `json:"type"`
	User string `json:"user"`
	Body any    `json:"-"`
}

type chatMessage struct {
	Message string `json:"message"`
}

type errorMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type commandMessage struct {
	Message string `json:"message"`
}
