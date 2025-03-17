package domain

import "github.com/gorilla/websocket"

type Player struct {
	Name string
	Conn *websocket.Conn

	send chan []byte
}
