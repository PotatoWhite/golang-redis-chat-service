package chat

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ChatRoom struct {
	Name    string
	Clients map[*websocket.Conn]struct{}
	lock    sync.Mutex
}

func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		Name:    name,
		Clients: make(map[*websocket.Conn]struct{}),
		lock:    sync.Mutex{},
	}
}

func (cr *ChatRoom) AddClient(conn *websocket.Conn) {
	cr.lock.Lock()
	cr.Clients[conn] = struct{}{}
	cr.lock.Unlock()
}

func (cr *ChatRoom) RemoveClient(conn *websocket.Conn) {
	cr.lock.Lock()
	delete(cr.Clients, conn)
	cr.lock.Unlock()
}

func (cr *ChatRoom) Broadcast(messageType int, msg []byte) {
	clientsCopy := make([]*websocket.Conn, 0, len(cr.Clients))
	for client := range cr.Clients {
		clientsCopy = append(clientsCopy, client)
	}

	for client := range clientsCopy {
		clientsCopy[client].WriteMessage(messageType, msg)
	}
}
