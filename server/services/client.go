package services

import (
	"log"
	"private-chat/core"
	"private-chat/events"
	"time"

	"github.com/gorilla/websocket"
)

// Client is the middleman between the websocket and the hub
type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send chan core.EventPayload
	UserId string
	Username string 
}

const (
	// Maximum message size allowed from peer 
	maxMessageSize = 1024  
	// Time allowed to read the message from the peer 
	readTimeout = time.Second * 60
)

func NewClientService() *Client{
	return &Client{}
}

// Pumps messages from the websocket to the hub.
// Hub will only do one thing and that is register the user to the hub
// This application ensures that there is at most one reader per connection 
// running as a goroutine.
func (c *Client) ReadPump() {
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
			c.newUserHandler(payload.EventPayload.(core.NewUserPayload))
		case events.DIRECT_MESSAGE:
			c.directMessageHandler(payload.EventPayload.(core.DirectMessagePayload))
		case events.DISCONNECT:
			disconnectHandler(payload.EventPayload.(core.DisconnectPayload))
		}
	} // end of for loop 
	
}

func (c *Client) newUserHandler(newUserPayload core.NewUserPayload) {
	// TODO Check if the user is logged in and if not don't do anything (just logged out user tried to create a conn)
	// Register the client 
	// Broadcast the connected users with the new user who has joined with the payload  
	log.Println("The new user has joined w/ username = ", newUserPayload.Username)
	// For new user send the chat list of all online users (except the user)
	var onlineUsers []core.NewUserPayload = []core.NewUserPayload{}
	for c := range c.Hub.Clients {
		onlineUsers = append(onlineUsers, core.NewUserPayload{Username: c.Username, UserId: c.UserId})
	}

	// Response sent to all the users except the joined user  
	response := core.EventPayload{
		EventName: events.NEW_USER,
		EventPayload: onlineUsers,
	}

	for client := range c.Hub.Clients {
		select {
		case client.Send <- response:
		default:
			close(client.Send)
			delete(c.Hub.Clients, client)
		}
	}
}

func (c *Client) directMessageHandler(directMessagePayload core.DirectMessagePayload) {
	// Extract out the UserId from payload to which the message needs to be sent 
	receiver := directMessagePayload.Receiver
	response := core.DirectMessageResponse{
		Sender: c.Username,
		Message: directMessagePayload.Message,
		Time: time.Now().String(),
	}
	// Loop over the hub clients and send the message to the specific user 	
	for c := range c.Hub.Clients {
		if c.UserId == receiver {
			c.Send <- core.EventPayload{EventName: events.DIRECT_MESSAGE, EventPayload: response}
			break
		}
	}
	log.Printf("There is a direct message for %v by %v", directMessagePayload.Receiver, directMessagePayload.Receiver)
}
func disconnectHandler(payload core.DisconnectPayload) {
	log.Println("The user : %v has disconnected", payload.Username)
}

// Pumps message from the hub to the websocket connection  
func (c *Client) writePump() {}