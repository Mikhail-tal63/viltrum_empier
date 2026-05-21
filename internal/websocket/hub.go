package websocket

import (
	"log"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]struct{}),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.addClient(c)
		case c := <-h.unregister:
			h.removeClient(c)
		}
	}
}

func (h *Hub) addClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	room, ok := h.rooms[c.WorkspaceID]
	if !ok {
		room = make(map[*Client]struct{})
		h.rooms[c.WorkspaceID] = room
	}
	room[c] = struct{}{}
}

func (h *Hub) removeClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	room, ok := h.rooms[c.WorkspaceID]
	if !ok {
		return
	}
	if _, exists := room[c]; !exists {
		return
	}
	delete(room, c)
	c.closeSend()
	if len(room) == 0 {
		delete(h.rooms, c.WorkspaceID)
	}
}

func (h *Hub) BroadcastToWorkspace(workspaceID string, payload []byte) {
	h.mu.RLock()
	room := h.rooms[workspaceID]
	clients := make([]*Client, 0, len(room))
	for c := range room {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		select {
		case c.send <- payload:
		default:
			log.Printf("ws: dropping slow client user=%s ws=%s", c.UserID, c.WorkspaceID)
			select {
			case h.unregister <- c:
			default:
			}
		}
	}
}
