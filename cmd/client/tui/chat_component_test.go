package tui_test

import (
	"testing"

	"github.com/dylanmccormick/ws-chat/cmd/client/tui"
)

func TestNewChatComponent(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *tui.ChatComponent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tui.NewChatComponent()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewChatComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}
