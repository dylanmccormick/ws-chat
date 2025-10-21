package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/dylanmccormick/ws-chat/cmd/server"
	"github.com/gorilla/websocket"
)

func Execute() {
	c := CreateConnection()
	defer c.Close()
	scanner := bufio.NewScanner(os.Stdin)
	bchan := make(chan []byte)
	currentRoom := "lobby"
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
		if len(input) > 0 && input[0] == '/' {
			tokens := strings.Split(input, " ")
			switch tokens[0] {
			case "/quit":
				break
			case "/create":
				msg := createCreateRoomMessage(tokens[1])
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Printf("We got an error writing: %s", err)
					panic(err)
				}
				continue
			case "/join":
				msg := createJoinRoomMessage(tokens[1])
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Printf("We got an error writing: %s", err)
					panic(err)
				}
				continue
			case "/switch":
				currentRoom = tokens[1]
				continue
			}
		}
		if input == "/quit" {
			break
		}

		if input == "" {
			continue
		}

		chat := createChatMessage(input, currentRoom)
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

func createJoinRoomMessage(name string) []byte {
	message := &server.Message{
		Typ: "command",
		Body: server.CommandMessage{
			Action: "JoinRoom",
			Target: name,
		},
	}
	msg, err := marshalJsonRepl(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func createCreateRoomMessage(name string) []byte {
	message := &server.Message{
		Typ: "command",
		Body: server.CommandMessage{
			Action: "CreateRoom",
			Target: name,
		},
	}
	msg, err := marshalJsonRepl(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func createChatMessage(input, room string) []byte {
	message := &server.Message{
		Typ: "chat",
		Body: server.ChatMessage{
			Message: input,
			Target:  room,
		},
	}
	msg, err := marshalJsonRepl(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func marshalJsonRepl(m *server.Message) ([]byte, error) {
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
