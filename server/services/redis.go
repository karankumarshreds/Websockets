package services

import (
	"private-chat/core"

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


func (r *RedisService) SaveMessageRedisToChat(message core.DirectMessagePayload) {
	/**
		-> create LSET for both the sender and the receiver 
		-> suppose the message is from A-> B 
		-> eg: LSET a_chats [B] and b_chats [A, C, D]
		-> create another LSET for the saving the chat data 
		-> eg: LSET "A.B" (the smaller uid will come first) [{},{},{}]
		-> while getting the all the chats for the user (say B)
		-> we LGET b_chats and create all combinations of chats:
		-> B.A, B.C, B.D (smaller uid coming before the ".") and once that is done we get
		-> all the chats using these possible keys 
		-> whichever of these combinations exists, return them with the latest message (last index)
	**/
}

func (r *RedisService) GetAllChats(forUserId string) {
	
}


func (r *RedisService) GetChatRedis( ) {

}
