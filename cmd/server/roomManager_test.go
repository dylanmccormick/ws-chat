package server

import (
	"reflect"
	"slices"
	"testing"
)

func TestRoomManager(t *testing.T) {
	rm := NewRoomManager()
	if reflect.TypeOf(rm) != reflect.TypeOf(&RoomManager{}) {
		t.Errorf("Unexpected type for RoomManager. expected=%T, got=%T", &RoomManager{}, rm)
	}

	tests := []struct {
		name string
	}{
		{"lobby"},
	}

	for _, tt := range tests {
		rm.AddRoom(tt.name)
		if rm.rooms[tt.name].Name != tt.name {
			t.Errorf("Unexpected name for room. expected=%s, got=%s", tt.name, rm.rooms[tt.name].Name)
		}

		rooms := rm.ListRooms()
		if !slices.Contains(rooms, tt.name) {
			t.Errorf("Room not added to room list. list=%v, expectedRoom=%s", rooms, tt.name)
		}

		r, err := rm.GetRoom(tt.name)
		if err != nil {
			t.Errorf("Unable to find room expected in room manager. expectedRoom=%s", tt.name)
		}
		if reflect.TypeOf(r) != reflect.TypeOf(&Room{}) {
			t.Errorf("Unexpected type for Room. expected=%T, got=%T", &Room{}, r)
		}

		err = rm.DeleteRoom(tt.name)
		if err != nil {
			t.Errorf("Unable to find room expected in room manager. expectedRoom=%s", tt.name)
		}
		_, ok := rm.rooms[tt.name]
		if ok {
			t.Errorf("Room was not successfully deleted from room manager. roomName=%s", tt.name)
		}

	}
}
