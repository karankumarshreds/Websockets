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
	// A -><- B 
	// B -><- C  
	// get all the conversations for b (with the latest message in each of the conversation)
	/**
	{
		B = [C, A] // b has chatted with C and A 
		// so the possible combinations to check from redis 
		-> B.C | C.B 
		-> B.A | A.B 
		A.B = [latest message]
		B.C = [latest message]
		A.C = [latest message] 
		-> get all the conversations with their latest message as preview for the user B 
	}
	**/
}

func (r *RedisService) GetAllChats(forUserId string) {
	
}


func (r *RedisService) GetChatRedis( ) {

}
