package services

import (
	"log"
	"private-chat/core"
	"private-chat/events"

	"github.com/go-redis/redis"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	rdb *redis.Client
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		Clients:    map[*Client]bool{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rdb: rdb,
	}
}

// Run acts like an interface between the readPump and writePump and updates Hub map
func (h *Hub) Run() {
	for { // infinite loop
		select {
		case client := <-h.Register:

			log.Println("Registering user with userid", client.UserId, "and username", client.Username)
			h.Clients[client] = true

			// create a list of online users (including the new user)
			var onlineUsers []core.NewUserPayload
			for c := range h.Clients {
				onlineUsers = append(onlineUsers, core.NewUserPayload{
					Username: c.Username,
					UserId: c.UserId,
				})
			}

			// all the users should be notified with the latest list of online users 
			for c := range h.Clients {
				// exclude broadcasting to the new user itself  
				if c.UserId != client.UserId {
					if len(FilterUser(onlineUsers, c.UserId)) > 0 {
						c.Send <- core.EventPayload{
							EventName: events.NEW_USER,
							// to make sure don't include userId of person to which this message will be sent 
							EventPayload: FilterUser(onlineUsers, c.UserId), 
						}
					}	
				} else { // for the newly joined user itself  
					if len(FilterUser(onlineUsers, c.UserId)) > 0 {
							c.Send <- core.EventPayload{
							EventName: events.NEW_USER,
							EventPayload: FilterUser(onlineUsers, c.UserId),
						}
					}

				}
			} 

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		}
	}
}


func FilterUser(users []core.NewUserPayload, userid string) []core.NewUserPayload {
	var filtered []core.NewUserPayload 
		for _, user := range users {
			if user.UserId != userid {
				filtered = append(filtered, user)
			}	
		}
		return filtered
}