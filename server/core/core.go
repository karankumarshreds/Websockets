package core

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Socket *websocket.Conn
	Send   string
	UserId string
}

type Hub struct {
	Clients 		map[*Client]bool
	Register 		chan *Client 
	Unregister 	chan *Client
}

type SocketEventStructure struct {
	EventName 		string `json:"eventName"`
	EventPayload 	string `json:"eventPayload"`
}
