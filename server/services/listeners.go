package services

import (
	"log"
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
			// var directMessagePayload core.DirectMessagePayload
			log.Printf("Received the message: %v for event: %v", msg.String(), msg.Channel)
		}
	}
}
