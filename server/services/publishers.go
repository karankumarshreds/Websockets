package services

import (
	"encoding/json"
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

// TODO you do not need to tell which user has left 
// the event is enough as the listeners will fetch the latest 
// list of the users from the redis map automatically anyways 
func (p *Publisher) NewUserPublisher(payload core.NewUserPayload) {
	log.Println("Publishing new user event via redis")
	if _payload, err := json.Marshal(payload); err != nil {
		log.Println("ERROR: Unable to marshall before publishing new user payload")
		return 
	} else {
		if err := p.rdb.Publish(string(events.NEW_USER), _payload).Err(); err != nil {
			log.Println("ERROR: Unable to publish NEW_USER event for payload", payload, err)
			return 
		} else {
			log.Println("Event successfully published for new user")
		}
	}
	
}

func (p *Publisher) UserDeletedPublisher() {
	log.Println("Publishing user deleted publisher via redis")
	if err := p.rdb.Publish(string(events.DELETED_USER), nil).Err(); err != nil {
		log.Println("ERROR: Unable to DELETED_USER event", err)
		return
	} else {
		log.Println("Event successfully published for deleted user")
	}
}

