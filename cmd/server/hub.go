package server

// The hub is the central event handler
type Hub struct {
	clients map[*User]bool

	messages chan []byte // all inbound messages for the hub. Will have user messages, commands, and announcements
}

// This is the event loop. All messages will come through the hub
func (h *Hub) run() {
	for {
		select{
		}

	// Process user messages

	// Process commands

	// Process announcements
}
}

