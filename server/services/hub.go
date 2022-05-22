package services

import (
	"encoding/json"
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
	log.Println("Creating a readpump for the new user")
	/* start listening for external messages */
	l := NewListeners(h.rdb, h)
	go l.NewUserListener()
	go l.DirectMessageListener()
	p := NewPublishers(h.rdb)
	
	for { // infinite loop
		select {
		case client := <-h.Register:



			/**
				## user joins to a serverA via hub
				-> create a new list of online users 
				-> broadcast all the local users 
				-> publish a message with new user event and payload 
			**/


			/* register the user in the local hub */ 
			log.Println("Hub.Run(): Registering user with userid", client.UserId, "and username", client.Username)
			h.Clients[client] = true
			
			h.BroadcastLocalUsers(client.UserId)	
			newUserPayload := core.NewUserPayload{UserId: client.UserId,Username: client.Username}
			h.UpdateRedisOnlineUsers(newUserPayload)		
			p.NewUserPublisher(newUserPayload)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		}
	} // end of infinite loop 
}

func (h *Hub) UpdateRedisOnlineUsers(newuser core.NewUserPayload) {
	log.Println("Updating the redis list of online users")
	var onlineUsers []core.NewUserPayload = []core.NewUserPayload{}
	for c := range h.Clients {
		onlineUsers = append(onlineUsers, core.NewUserPayload{
			Username: c.Username,
			UserId: c.UserId,
		})
	}
	onlineUsers = append(
		onlineUsers, 
		core.NewUserPayload{Username: newuser.Username, UserId: newuser.UserId},
	)
	if onlineUsersMarshalled, marshalErr := json.Marshal(onlineUsers); marshalErr != nil {
		log.Println("ERROR: Unable to marhshal list of online users", marshalErr) 
		return 
	} else {
		if rdbErr := h.rdb.Set(REDIS_KEYS.online_users, onlineUsersMarshalled, 0).Err(); rdbErr != nil {
			log.Println("ERROR: Unable to save online users list on redis", rdbErr)
			return
		}
	}
}

// broadcast all the local users in the server memory about the new user 
func (h *Hub) BroadcastLocalUsers(userid string) {
	// create a list of online users (including the new user)
	log.Println("Hub.Run(): Creating a list of online users")
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
		if c.UserId != userid {
			if len(h.FilterUser(onlineUsers, c.UserId)) > 0 {
				log.Println("Hub.Run(): Emitting online users to everyone except the new user")
				c.Send <- core.EventPayload{
					EventName: events.NEW_USER,
					// to make sure don't include userId of person to which this message will be sent 
					EventPayload: h.FilterUser(onlineUsers, c.UserId), 
				}
			} else {
				log.Println("Hub.Run(): No other users to send online users list")
			}
		} else { // for the newly joined user itself  
			if len(h.FilterUser(onlineUsers, c.UserId)) > 0 {
				log.Println("Hub.Run(): Emitting the list of online users to new user")
				c.Send <- core.EventPayload{	
					EventName: events.NEW_USER,
					EventPayload: h.FilterUser(onlineUsers, c.UserId),
				}
			} else {
				c.Send <- core.EventPayload{
					EventName: events.NEW_USER,
					EventPayload: []interface{}{},
				}
			}
		}
	} 
}


func (h *Hub) FilterUser(users []core.NewUserPayload, userid string) []core.NewUserPayload {
	var filtered []core.NewUserPayload 
		for _, user := range users {
			if user.UserId != userid {
				filtered = append(filtered, user)
			}	
		}
		return filtered
}