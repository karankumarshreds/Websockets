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

				// make sure the user is not local (make sure the event was not published by us) 
				isLocalUser := false 
				log.Println("Checking for possibility of a local user")
				for c := range l.hub.Clients {
					if c.UserId == newUserPayload.UserId {
						isLocalUser = true
						break
					}

				}
				if isLocalUser {
					log.Println("The user is local, redis event publish not required")
					return
				}
				// as soon as we get an event of new user, we fetch the latest list of the online users from redis 
				log.Println("A new user joined another server, getting new list of online users from redis")
				if data, rdbErr := l.rdb.Get(REDIS_KEYS.online_users).Result(); rdbErr != nil {
					log.Println("ERROR: Unable to get list of the online users from redis", rdbErr)
					return
				} else {
					var onlineUsers []core.NewUserPayload
					if err := json.Unmarshal([]byte(data), &onlineUsers); err != nil {
						log.Println("Unable to Unmarshal the list of online users")
						return 
					} 
					// socket send all the local users with the updated list of users 
					for c := range l.hub.Clients {
						c.Send <- core.EventPayload{
							EventName: events.NEW_USER,
							// to make sure don't include userId of person to which this message will be sent 
							EventPayload: l.hub.FilterUser(onlineUsers, c.UserId), 
						}
					}
				}
				return // end of new logic 
			}
		}
	}
}
