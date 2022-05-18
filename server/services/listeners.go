package services

import (
	"encoding/json"
	"log"
	"private-chat/core"
	"private-chat/events"

	"github.com/go-redis/redis"
)

type Listener struct {
	rdb *redis.Client
	hub *Hub
}

func NewListeners(rdb *redis.Client, hub *Hub) *Listener {
	return &Listener{rdb, hub}
}

func (l *Listener) DirectMessageListener() {
	log.Println("Listeners.go := Listening for direct messages...")
	s := l.rdb.Subscribe(string(events.DIRECT_MESSAGE))
	for {
		if msg, err := s.ReceiveMessage(); err != nil {
			log.Println("ERROR: Error receiving message for the event name", events.DIRECT_MESSAGE, err)
		} else {
			// var directMessagePayload core.DirectMessagePayload
			log.Printf("Received the message: %v for event: %v", msg.String(), msg.Channel)
		}
	}
}

func (l *Listener) NewUserListener() {
	log.Println("Listeners.go := Listening for new user...")
	s := l.rdb.Subscribe(string(events.NEW_USER))
	for {
		if msg, err := s.ReceiveMessage(); err != nil {
			log.Println("ERROR: Error receiving message for the event name", events.DIRECT_MESSAGE, err)
		} else {
			log.Printf("Received the message: %v for event: %v", msg.String(), msg.Channel)
			log.Println("Unmarshalling received message from redis listener")
			var newUserPayload core.NewUserPayload
			if unmarshalErr := json.Unmarshal([]byte(msg.Payload), &newUserPayload); unmarshalErr != nil {
				log.Println("ERROR: Unable to unmarshal received message from redis listener", unmarshalErr)
				return 
			} else {
				isLocalUser := true 
				for c := range l.hub.Clients {
					if c.UserId == newUserPayload.UserId {
						isLocalUser = false
						break
					}

				}
				if isLocalUser {
					// create a new client using the incoming payload 
					// we don't have to register this user as it is already register with some other server 
					var onlineUsers []core.NewUserPayload = []core.NewUserPayload{}
					for c := range l.hub.Clients {
						onlineUsers = append(onlineUsers, core.NewUserPayload{
							Username: c.Username,
							UserId: c.UserId,
						})
					}
					onlineUsers = append(
						onlineUsers, 
						core.NewUserPayload{Username: newUserPayload.Username, UserId: newUserPayload.UserId},
					)
					log.Println("Broadcasting to local users after receiving NEW_USER event to", onlineUsers)
					for c := range l.hub.Clients {
						c.Send <- core.EventPayload{
							EventName: events.NEW_USER,
							// to make sure don't include userId of person to which this message will be sent 
							EventPayload: l.hub.FilterUser(onlineUsers, c.UserId), 
						}
					}
					
				} else {
					log.Println("Broadcast not required as the user is local")
				}
				
				
			}
		}
	}
}
