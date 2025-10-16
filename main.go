package main

import (
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

func main() {
	fmt.Println("Hello and welcome to websocket chat")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", "error", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("An error occurred upgrading the http connection", "error", err)
		panic(err)
	}

	go timedWrite(conn)

}

func timedWrite(conn *websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			ws, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("An error occurred with NextWriter: ", "error", err)
			}

			ws.Write([]byte("Test"))
			ws.Close()
		default:
			continue
		}
	}
}
