package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
	"strings"

	prot "github.com/dylanmccormick/ws-chat/internal/protocol"
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
				msg := CreateCreateRoomMessage(tokens[1])
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Printf("We got an error writing: %s", err)
					panic(err)
				}
				continue
			case "/join":
				msg := CreateJoinRoomMessage(tokens[1])
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Printf("We got an error writing: %s", err)
					panic(err)
				}
				continue
			case "/switch":
				currentRoom = tokens[1]
				continue
			case "/list":
				msg := CreateListRoomMessage()
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Printf("We got an error writing: %s", err)
					panic(err)
				}
			}
		}
		if input == "/quit" {
			break
		}

		if input == "" {
			continue
		}

		chat := CreateChatMessage(input, currentRoom)
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
	userNumber := fmt.Sprintf("%6d", rand.IntN(999999))
	c.WriteMessage(websocket.TextMessage, []byte("TestUser"+userNumber))
	return c
}

func CreateJoinRoomMessage(name string) []byte {
	message := &prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Action: "JoinRoom",
			Target: name,
		},
	}
	msg, err := MarshalJson(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func CreateGetUsersMessage(room string) []byte {
	message := &prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Action: "ListRoomUsers",
			Target: room,
		},
	}
	msg, err := MarshalJson(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func CreateListRoomMessage() []byte {
	message := &prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Action: "ListMyRooms",
		},
	}
	msg, err := MarshalJson(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func CreateCreateRoomMessage(name string) []byte {
	message := &prot.Message{
		Typ: "command",
		Body: prot.CommandMessage{
			Action: "CreateRoom",
			Target: name,
		},
	}
	msg, err := MarshalJson(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func CreateChatMessage(input, room string) []byte {
	message := &prot.Message{
		Typ: "chat",
		Body: prot.ChatMessage{
			Message: input,
			Target:  room,
		},
	}
	msg, err := MarshalJson(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func MarshalJson(m *prot.Message) ([]byte, error) {
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
