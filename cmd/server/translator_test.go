package server

import (
	"context"
	"reflect"
	"testing"
)

func TestUnmarshalMessage(t *testing.T) {
	tests := []struct {
		input              string
		expectedTypeString string
		expectedType       any
	}{
		{`{"type": "chat", "body": {"test": "ok"}}`, "chat", ChatMessage{}},
		{`{"type": "error", "body": {"test": "ok"}}`, "error", ErrorMessage{}},
		{`{"type": "command", "body": {"test": "ok"}}`, "command", CommandMessage{}},
	}

	for _, tt := range tests {
		var m Message
		err := m.UnmarshalJSON([]byte(tt.input))
		if err != nil {
			t.Errorf("error occurred when unmarshalling json: %s", err)
		}
		if m.Typ != tt.expectedTypeString {
			t.Errorf("Received unexpected type from unmarshal. got=%s expected=%s", m.Typ, tt.expectedTypeString)
		}
		if reflect.TypeOf(m.Body) != reflect.TypeOf(tt.expectedType) {
			t.Errorf("Received unexpected type in message body. got=%s expected=%s", reflect.TypeOf(m.Body), reflect.TypeOf(tt.expectedType))
		}
	}
}

func TestMarshalMessage(t *testing.T) {
	translator := Translator{}
	tests := []struct {
		input        Message
		expectedData string
	}{
		{Message{Typ: "chat", Body: ChatMessage{Message: "test", Target: "lobby"}}, `{"type":"chat","body":{"message":"test","target":"lobby"}}`},
	}

	for _, tt := range tests {
		dat, err := translator.MessageToBytes(context.TODO(), tt.input)
		if err != nil {
			t.Errorf("error occurred marshalling json: %s", err)
		}

		if string(dat) != tt.expectedData {
			t.Errorf("Json output does not match expected string: expected=%s got=%s", tt.expectedData, string(dat))
		}
	}
}
