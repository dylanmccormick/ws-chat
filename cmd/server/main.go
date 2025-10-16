package server

import (
	"bytes"
	"fmt"
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
	send        chan string
}

type Room struct {
	Name        string
	Users       []User
	MessageChan chan string
	// Should there be a message chan??? ... probably
}

type Server struct {
	Rooms []Room
}

func StartServer() {
	s := Server{[]Room{}}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(w, r)
	})
	room := createRoom()
	s.Rooms = append(s.Rooms, room)
	// go s.updateRooms()
	go s.Rooms[0].updateRoom()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", "error", err)
	}
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	slog.Info("upgrading the server")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("An error occurred upgrading the http connection", "error", err)
		panic(err)
	}

	slog.Info("Creating a new user")
	user, err := s.createUser(conn)
	s.Rooms[0].Users = append(s.Rooms[0].Users, user)

	slog.Info("Created user: ", "username", user.username)

	go timedWrite(user)
	go reader(user)
}

func (s *Server) createUser(conn *websocket.Conn) (User, error) {
	ws, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error("An error occurred with NextWriter: ", "error", err)
	}

	ws.Write([]byte("Please send your username"))
	ws.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		slog.Error("That username is bogus", "error", err)
		return User{}, err
	}

	username := bytes.TrimSpace(bytes.ReplaceAll(message, []byte("\n"), []byte(" ")))
	ws, err = conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error("An error occurred with NextWriter: ", "error", err)
	}
	fmt.Fprintf(ws, "Welcome to the lobby, %s", username)
	ws.Close()


	return User{username: string(username), conn: conn, currentRoom: s.Rooms[0], send: make(chan string)}, nil
}

func reader(u User) {
	slog.Info("Starting reader")
	for {
		_, message, err := u.conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", "error", err)
			break
		}

		// Here we will parse the command (if the message starts with /)
		// If not, that is a message for the chat
		message = bytes.TrimSpace(bytes.ReplaceAll(message, []byte("\n"), []byte(" ")))

		// Here we should send this to some sort of message broker. It handles incoming and outgoing chats for each client
		// That way it can reference the room which it is in and handle messages from that room. Or announcements from the server
		slog.Info("Got a message", "message", message)
		u.currentRoom.MessageChan <- fmt.Sprintf("%s: %s", u.username, string(message))
		slog.Info("Message sent to channel")
	}
}

func createRoom() Room {
	return Room{
		Name:        "lobby",
		Users:       []User{},
		MessageChan: make(chan string),
	}
}

func timedWrite(u User) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case msg := <-u.send:
			slog.Info("Got a message from u.send channel", "user", u.username)
			ws, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("An error occurred with NextWriter: ", "error", err)
			}

			ws.Write([]byte(msg))
			ws.Close()
		default:
			continue
		}
	}
}

func (r *Room) updateRoom() {
	for {
		select {
		case msg := <-r.MessageChan:
			slog.Info("Got a message in update room channel", "room", r.Name, "message", msg)
			for _, u := range r.Users {
				u.send <- msg
			}
		default:
			continue
		}
	}
}

func (s *Server) updateRooms() {
	// This is going to be a problem. We don't want to spin off a goroutine for each room each loop. Infinite goroutines == bad
}
