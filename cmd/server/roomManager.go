package server

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type Room struct {
	Name  string
	Users []*User
}

type RoomManager struct {
	rooms map[string]*Room
	mux   sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
		mux:   sync.Mutex{},
	}
}

func (r *RoomManager) ListRooms() []string {
	r.mux.Lock()
	defer r.mux.Unlock()
	// I don't think these really need to be sorted, but maybe?
	return slices.Sorted(maps.Keys(r.rooms))
}

func (r *RoomManager) AddRoom(name string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.rooms[name]; ok {
		return fmt.Errorf("the room %s already exists", name)
	}
	r.rooms[name] = &Room{Name: name, Users: []*User{}}
	return nil
}

func (r *RoomManager) DeleteRoom(name string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.rooms[name]; !ok {
		return fmt.Errorf("the room %s does not exist", name)
	}
	delete(r.rooms, name)
	return nil
}

func (r *RoomManager) GetRoom(name string) (*Room, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	if room, ok := r.rooms[name]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("the room %s does not exist", name)
}
