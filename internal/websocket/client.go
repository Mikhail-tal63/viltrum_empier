package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	UserID string

	WorkspaceID string
	GroupChatID string

	closeOnce sync.Once
}


func (c *Client) closeSend() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}
