package services

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

/************************************
 -> user data
 {
	 [userid]: {
			username: string,
			socket: websocket,
	 }
 }

 -> chat data
 {
	 [receiver.sender]: {
		 [chunk1]: [ ...10 messages ]
	 }

************************************/

type RedisService struct {
	rdb *redis.Client
}

type UserNode struct {
		Username string
		Socket *websocket.Conn
}

func NewRedisService(rdb *redis.Client) *RedisService {
	return &RedisService{rdb}
}


func (r *RedisService) SaveMessageRedis() {

}

func (r *RedisService) GetAllChats(forUserId string) {
	
}


func (r *RedisService) GetChatRedis( ) {

}
