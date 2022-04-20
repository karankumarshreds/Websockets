package service

import (
	"log"
	"private-chat/core"
	"private-chat/events"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send string
	UserId string
}

const (
	// Maximum message size allowed from peer 
	maxMessageSize = 1024  
	// Time allowed to read the message from the peer 
	readTimeout = time.Second * 60
)

// Pumps messages from the websocket to the hub.
// This application ensures that there is at most one reader per connection 
// running as a goroutine.
func (c *Client) readPump() {
	defer func() {
		// unregister client on while terminating the goroutine by sending the client to unregister channel 
		c.Hub.Unregister <- c 
		// close the connection 
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(readTimeout)) // message will not be read 60 seconds after recieving 
	// send the message to the client (ping) to get a response (pong) and update the deadline if the pong is recieved 
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(readTimeout)); return nil })
	
	for { // infinite for loop
		// listen for any incoming message from the websocket connection 
		var payload core.EventPayload
		if err := c.Conn.ReadJSON(&payload); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected error : ", err)
				break
			}
		}
		switch payload.EventName {
		case events.NEW_USER:
			newUserHandler(payload.EventPayload.(core.NewUserPayload))
		case events.DIRECT_MESSAGE:
			directMessageHandler(payload.EventPayload.(core.DirectMessagePayload))
		}
	} // end of for loop 
	
}

func newUserHandler(payload core.NewUserPayload) {
	log.Println("The new user has joined w/ username = ", payload.Username)
}
func directMessageHandler(payload core.DirectMessagePayload) {
	log.Printf("There is a direct message for %v by %v", payload.Receiver, payload.Receiver)
}


func (c *Client) writePump() {}