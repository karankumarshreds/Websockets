package services

import (
	"log"
	"private-chat/core"
	"private-chat/events"

	"github.com/go-redis/redis"
)

type Publisher struct {
	rdb *redis.Client
}

func NewPublishers(rdb *redis.Client) *Publisher {
	return &Publisher{rdb}
}

func (p *Publisher) NewUserPublisher(payload core.NewUserPayload) {
	log.Println("Publishing new user event via redis")
	if err := p.rdb.Publish(string(events.NEW_USER), payload).Err(); err != nil {
		log.Println("ERROR: Unable to publish NEW_USER event for payload", payload, err)
		return 
	}
	log.Println("Event successfully published for new user")
}

