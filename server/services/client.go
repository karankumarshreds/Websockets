package services

import (
	"encoding/json"
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
	// rdb *redis.Client
	redisService *RedisService
}


const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Maximum message size allowed from peer 
	maxMessageSize = 1024  
	// Time allowed to read the message from the peer 
	readTimeout = time.Second * 10
	// Send pings to peer with this period. Must be less than pongWait (taking 90% of readTimeout)
	pingPeriod = (readTimeout * 9) / 10
)

func NewClientService(
	Hub *Hub,
	Conn *websocket.Conn,
	Send chan core.EventPayload,
	UserId string,
	Username string ,
	redisService *RedisService,
) *Client{
	return &Client{ Hub,Conn,Send,UserId,Username,redisService}
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
	c.Conn.SetReadDeadline(time.Now().Add(readTimeout)) // message will not be read 60 seconds after recieving/processing 
	// send the message to the client (ping) to get a response (pong) and update the deadline if the pong is recieved 
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(readTimeout)); return nil })
	
	// listen for any incoming message from the websocket connection 
	for { // infinite for loop
		var payload core.EventPayload
		if err := c.Conn.ReadJSON(&payload); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected error : ", err)
				break
			} else {
				log.Println("Error connection broken by client", err)
				break
			}
		}
		switch payload.EventName {
		case events.DIRECT_MESSAGE:
			var directMessagePayload core.DirectMessagePayload
			c.unmarshalPayload(payload.EventPayload, &directMessagePayload)
			if !c.MessageBroadcastRequired(directMessagePayload.UserId, string(events.DIRECT_MESSAGE), directMessagePayload) {
				c.directMessageHandler(directMessagePayload)
			}
		case events.DISCONNECT:
			var disconnectPayload core.DisconnectPayload
			c.unmarshalPayload(payload.EventPayload, &disconnectPayload)
			if !c.MessageBroadcastRequired(disconnectPayload.UserId, string(events.DISCONNECT), disconnectPayload) {
				c.disconnectHandler(disconnectPayload)
			}
		}
	} // end of for loop 
}

// userid => id of user which is the receiver of the event  
func (c *Client) MessageBroadcastRequired(userid string, eventName string, eventPayload interface{}) bool {
	log.Println("Checking if the broadcast is required or not")
	if eventName == string(events.NEW_USER) {
		return false 
	}
	var uid string 
	if eventName == string(events.DIRECT_MESSAGE) {
		log.Printf("BroadcastCheck => eventName : %v", eventName)
		uid = eventPayload.(core.DirectMessagePayload).UserId
	}
	if eventName == string(events.DISCONNECT) {
		log.Printf("BroadcastCheck => eventName : %v", eventName)
		uid = eventPayload.(core.DisconnectPayload).UserId
	}
	// checks if the client(receiver) is there in the local memory hub 
	for c := range c.Hub.Clients {
		if c.UserId == uid {
			return false 
		}
	} 
	log.Println("Publishing the event with event name", eventName)
	if data, marshalErr  := json.Marshal(eventPayload); marshalErr != nil {
		log.Println("ERROR: Marshal error before publishing to redis", marshalErr)
	} else {
		if err := c.redisService.rdb.Publish(eventName, data).Err(); err != nil {
		log.Println("ERROR: Unable to publish the event", eventName, err)
	}
	}
	return true 
}


/* handler function for the direct message type */
func (c *Client) directMessageHandler(directMessagePayload core.DirectMessagePayload) {
	// Extract out the UserId from payload to which the message needs to be sent 
	receiver := directMessagePayload.UserId
	// creating response for the receiver 
	response := core.DirectMessageResponse{
		UserId: directMessagePayload.UserId,
		Sender: c.Username,
		Message: directMessagePayload.Message,
		Time: time.Now().String(),
	}
	// Loop over the hub clients and send the message to the specific user 	
	for client := range c.Hub.Clients {
		if client.UserId == receiver {
			// save the chat using redis service 
			_message := SaveMessageArg{
				Sender: struct{Username string; UserId string}{
					c.getUsername(directMessagePayload.Sender),
					directMessagePayload.Sender,
				},
				Receiver: struct{Username string; UserId string}{
					c.getUsername(directMessagePayload.UserId),
					directMessagePayload.UserId,
				},
				Message: directMessagePayload.Message,
				Time: time.Now().String(),
			}
			c.redisService.SaveMessageRedisToChat(_message)
			client.Send <- core.EventPayload{EventName: events.DIRECT_MESSAGE, EventPayload: response}
			break
		}
	}
	log.Printf("There is a direct message for %v by %v", directMessagePayload.UserId, directMessagePayload.Sender)
}


/* handler function for the disconnect type */
func (c *Client) disconnectHandler(disconnectedUserPayload core.DisconnectPayload) {
	c.Hub.Unregister <- c
	// Broadcast all users with disconnected user list 
	for client := range c.Hub.Clients {
		client.Send <- core.EventPayload{EventName: events.DELETED_USER, EventPayload: disconnectedUserPayload.UserId}	
	}
	log.Printf("The user : %v has disconnected", disconnectedUserPayload.Username)
}

// Pumps message from the hub to the websocket connection  
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
			ticker.Stop()
			c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <- c.Send: 
		// log.Println("WritePump(): checking if the broadcast is required")
			// if c.MessageBroadcastRequired(string(message.EventName),string(message.EventName),message.EventPayload) {
			// 	return // return early if the broadcasting is required 
			// }
			// every message will have an eventName attached to it 
			log.Println("WritePump(): writing to the client", message.EventPayload)
			// Setting a deadline to write this message to the websocket 
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Assuming the hub closed the channel 
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			if err := c.Conn.WriteJSON(message); err == nil {
				log.Println("WritePump(): message write to client successful", message)	
			} else {
				log.Println("WritePump(): ERROR Could not write message to client => ", err)
				return 
			} 

			// if w, err := c.Conn.NextWriter(websocket.TextMessage); err != nil {
			// 	return 
			// } else {
			// 	reqBodyBytes := new(bytes.Buffer)
			// 	if encodeErr := json.NewEncoder(reqBodyBytes).Encode(message); encodeErr != nil {
			// 		return 
			// 	} else {
			// 		log.Println("Message writing to the client", message)
			// 		w.Write(reqBodyBytes.Bytes())
			// 	}
			// }
		case <- ticker.C:
			// Setting a deadline to write this message to the websocket 
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Time to send another ping message (for which also, we've put a deadline above)
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}	

func (c *Client) getUsername (userid string) string {
	var username string 
	for client := range c.Hub.Clients {
		if client.UserId == userid {
			username = client.Username
		}
	}
	return username
}

func (c *Client) unmarshalPayload(payload interface{}, v interface{}) {
	// Your interface (payload) value holds a map, you can't convert that to a struct. 
	// Use json.RawMessage when unmarshaling, and when you know what type you need, 
	// do a second unmarshaling into that type. You may put this logic into an 
	// UmarshalJSON() function to make it automatic. 
	if data, err := json.Marshal(payload); err != nil {
		log.Println("ERROR: Unable to marshal the incoming payload", err)
	} else {
		if err := json.Unmarshal(data, v); err != nil {
			log.Println("ERROR: Unable to unmarshal the incoming payload", err)
		}
	}	
}