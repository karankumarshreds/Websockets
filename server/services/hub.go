package services

import (
	"encoding/json"
	"errors"
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

			newUserPayload := core.NewUserPayload{UserId: client.UserId,Username: client.Username}
			onlineUsers := h.UpdateRedisOnlineUsers(newUserPayload)		
			if onlineUsers != nil {
				h.BroadcastLocalUsers(client.UserId, *onlineUsers)	
				p.NewUserPublisher(newUserPayload)
			}
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		}
	} // end of infinite loop 
}

func (h *Hub) UpdateRedisOnlineUsers(newuser core.NewUserPayload) *[]core.NewUserPayload {

	// TODO you have to merge the list of online users 
	// create a map of online users 
	log.Println("Creating a map of online users")
	var onlineUsersMap map[string]core.NewUserPayload 

	// get the already existing map and update it and then save it 
	onlineUsersRedisMap, redisErr := h.rdb.Get(REDIS_KEYS.online_users_map).Result(); 

	if errors.Is(redisErr, redis.Nil) {
		log.Println("Redis online users map not found, creating one")
		onlineUsersMap[newuser.UserId] = newuser;
	} else if redisErr != nil  {
		log.Println("Cannot load the online users map from redis", redisErr)
		return nil
	} else { 
		if err := json.Unmarshal([]byte(onlineUsersRedisMap), &onlineUsersMap); err != nil {
			log.Println("Unable to unmarshal the online users from redis", err)
			return nil
		}
		// if all good, merge the state with new user payload 
		onlineUsersMap[newuser.UserId] = newuser;
	}
	
	log.Println("Marshalling the map of online users")
	if onlineUsersMapMarshalled, marshalErr := json.Marshal(onlineUsersMap); marshalErr != nil {
		log.Println("ERROR: Unable to marshal the map of online users")
		return nil
	} else {
		if err := h.rdb.Set(REDIS_KEYS.online_users_map, onlineUsersMapMarshalled,0).Err(); err != nil {
			log.Println("Unable to create redis user map", err)
			return nil
		}
	}
	

	log.Println("Updating the redis list of online users")
	var onlineUsers []core.NewUserPayload = []core.NewUserPayload{}
	for _, userData := range onlineUsersMap {
		onlineUsers = append(onlineUsers, userData)
	}

	if onlineUsersMarshalled, marshalErr := json.Marshal(onlineUsers); marshalErr != nil {
		log.Println("ERROR: Unable to marhshal list of online users", marshalErr) 
		return nil
	} else {
		if rdbErr := h.rdb.Set(REDIS_KEYS.online_users, onlineUsersMarshalled, 0).Err(); rdbErr != nil {
			log.Println("ERROR: Unable to save online users list on redis", rdbErr)
			return nil
		}
	}
	return &onlineUsers
}


func (h *Hub) BroadcastLocalUsers(userid string, onlineUsers []core.NewUserPayload) {

	// create a list of online users (including the new user)
	log.Println("Hub.Run(): Creating a list of online users")
	// var onlineUsers []core.NewUserPayload
	// for c := range h.Clients {
	// 	onlineUsers = append(onlineUsers, core.NewUserPayload{
	// 		Username: c.Username,
	// 		UserId: c.UserId,
	// 	})
	// }

	// TODO you have to merge the list of online users 

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