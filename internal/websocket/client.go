package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	UserID string
}
