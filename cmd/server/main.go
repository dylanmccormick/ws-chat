package server

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type User struct {
	conn        *websocket.Conn
	username    string
	currentRoom Room
	send        chan []byte
}

type Server struct {
	Rooms []Room
	Hub   *Hub
}

func StartServer() {
	slog.Info("Starting server")
	hub := NewHub()
	s := Server{[]Room{}, hub}
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", "error", err)
	}
}

// TODO: This is where context should be created and passed around
func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	slog.Info("upgrading the server")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("An error occurred upgrading the http connection", "error", err)
		panic(err)
	}

	slog.Info("Creating a new user")
	user, err := s.createUser(conn)
	// TODO: This needs to be done differently! I think a userManager or something
	// Or send a newUser message to the hub and it'll handle user-specific actions!

	slog.Info("Created user: ", "username", user.username)

	go timedWrite(user)
	// go reader(user, s.Hub.messages)
	s.Hub.register <- user
}

func (s *Server) createUser(conn *websocket.Conn) (*User, error) {
	// We're instead going to create an anonymous user and send them through the registration channel (before they can get messages)
	u := &User{
		conn: conn,
		send: make(chan []byte, 10),
	}
	return u, nil
}

func reader(u *User, highway chan InternalMessage) {
	slog.Info("Starting reader")
	t := Translator{}
	for {
		_, data, err := u.conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", "error", err)
			break
		}
		data = bytes.TrimSpace(bytes.ReplaceAll(data, []byte("\n"), []byte(" ")))
		slog.Info("Got a message", "message", data)
		message, err := t.BytesToMessage(context.TODO(), data)
		if err != nil {
			slog.Error("Error turning data ([]bytes) into Message", "data", string(data), "location", "reader")
		}
		message.EnrichWithUser(u)
		highway <- message
		slog.Info("Message sent to channel")
	}
}

func timedWrite(u *User) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case data := <-u.send:
			slog.Info("Got a message from u.send channel", "user", u.username)
			WriteToConn(u.conn, data)
		default:
			continue
		}
	}
}
