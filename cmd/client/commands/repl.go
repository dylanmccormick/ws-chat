package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"

	"github.com/dylanmccormick/ws-chat/cmd/server"
	"github.com/gorilla/websocket"
)

func Execute() {
	c := CreateConnection()
	defer c.Close()
	scanner := bufio.NewScanner(os.Stdin)
	bchan := make(chan []byte)
	go readAndPrint(c, bchan)
	for {
		for {
			b := false
			select {
			case message := <-bchan:
				fmt.Println(string(message))
			default:
				b = true
			}
			if b {
				break
			}
		}
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "/quit" {
			break
		}

		chat := createChatMessage(input)
		err := c.WriteMessage(websocket.TextMessage, chat)
		if err != nil {
			fmt.Printf("We got an error writing: %s", err)
			panic(err)
		}
	}
}

func readAndPrint(c *websocket.Conn, bchan chan []byte) {
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}
		data = bytes.TrimSpace(bytes.ReplaceAll(data, []byte("\n"), []byte(" ")))
		bchan <- data
	}
}

func CreateConnection() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	c.WriteMessage(websocket.TextMessage, []byte("TestUser"))
	return c
}

func createChatMessage(input string) []byte {
	message := &server.Message{
		Typ: "chat",
		Body: &server.ChatMessage{
			Message: input,
			Target:  "lobby",
		},
	}
	msg, err := message.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return msg
}
