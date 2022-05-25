package services

import (
	"encoding/json"
	"log"
	"private-chat/core"
	"private-chat/events"
	"private-chat/utils"

	"github.com/go-redis/redis"
)

type Listener struct {
	rdb *redis.Client
	hub *Hub
	redisService *RedisService
}

func NewListeners(rdb *redis.Client, hub *Hub, redisService *RedisService) *Listener {
	return &Listener{rdb, hub, redisService}
}

func (l *Listener) DirectMessageListener() {
	log.Println("Listeners.go := Listening for direct messages...")
	s := l.rdb.Subscribe(string(events.DIRECT_MESSAGE))
	for {
		if msg, err := s.ReceiveMessage(); err != nil {
			log.Println("ERROR: Error receiving message for the event name", events.DIRECT_MESSAGE, err)
		} else {
			utils.CustomLogger("Received the message: %v for event: %v", msg.String(), msg.Channel)
			var directMessagePayload core.DirectMessagePayload
			// check if the receiver is in the list of online users 
			// if yes then send the message to the receiver 
			if err := json.Unmarshal([]byte(msg.Payload), &directMessagePayload); err != nil {
				utils.CustomLogger("Unable to unmarshal the DIRECT_MESSAGE payload", err)
				return
			}
			utils.CustomLogger("Unmarshal successful, sending the message the receiver", directMessagePayload.UserId)
			for client := range l.hub.Clients {
				if client.UserId == directMessagePayload.UserId {
					if err := client.Conn.WriteJSON(core.EventPayload{
						EventName: events.DIRECT_MESSAGE,
						EventPayload: directMessagePayload,
					}); err != nil {
						utils.CustomLogger("ERROR: Unable to send message to the receiver", err)
						return
					} else {
						utils.CustomLogger("Message successfully sent to", client.UserId)
					}
				}
			}
		}
	}
}


func (l *Listener) NewUserListener() {
	log.Println("Listeners.go := Listening for new user...")
	s := l.rdb.Subscribe(string(events.NEW_USER))

	for { // start of infinite loop
		msg, err := s.ReceiveMessage(); 
		if err != nil {
			log.Println("ERROR: Error receiving message for the event name", events.DIRECT_MESSAGE, err)
			return
		} 
		log.Printf("Received the message: %v for event: %v", msg.String(), msg.Channel)
		log.Println("Unmarshalling received message from redis listener")
		var newUserPayload core.NewUserPayload
		if unmarshalErr := json.Unmarshal([]byte(msg.Payload), &newUserPayload); unmarshalErr != nil {
			log.Println("ERROR: Unable to unmarshal received message from redis listener", unmarshalErr) 
			return
		} 

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
		} else {
			// as soon as we get an event of new user, we fetch the latest list of the online users from redis 
			utils.CustomLogger("A new user joined another server, getting new list of online users from redis")
			onlineUsersRedisMap, err := l.redisService.GetOnlineUsersRedisMap(); 
			if err != nil {
				utils.CustomLogger("Unable to listen to NEW_USER event")  
			}
			onlineUsers :=  []core.NewUserPayload{}
			for _, payload := range onlineUsersRedisMap {
				onlineUsers = append(onlineUsers, payload)
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
	} // end of infinite loop
}
