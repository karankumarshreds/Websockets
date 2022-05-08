package services

import (
	"encoding/json"
	"fmt"
	"log"
	"private-chat/core"
	"private-chat/utils"
	"strconv"

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


/**
	-> check if the LSET for A or B exists or not 
	-> if not exists then 
		-> create LSET for both the sender and the receiver 
		-> suppose the message is from A-> B 
		-> eg: LSET a_chats [B] and b_chats [A, C, D]
	-> else
	  -> check if the userid in the array exists or not 
		-> if not then update the array  
	-> create another LSET for the saving the chat data 
	-> eg: LSET "A.B" (the smaller uid will come first) [{},{},{}]
	-> while getting the all the chats for the user (say B)
	-> we LGET b_chats and create all combinations of chats:
	-> B.A, B.C, B.D (smaller uid coming before the ".") and once that is done we get
	-> all the chats using these possible keys 
	-> whichever of these combinations exists, return them with the latest message (last index)
**/

func (r *RedisService) SaveMessageRedisToChat(message core.DirectMessagePayload) {

}

func (r *RedisService) GetAllChats(forUserId string) {
	
}


func (r *RedisService) GetChatRedis() {

}

func (r *RedisService) CreateChatCombinations(forUserId string, withUserId string)  {
	c := fmt.Sprintf("%v_CHATS", forUserId)
	if list, err  := r.rdb.LRange(c, 0, -1).Result(); err != nil {
		log.Println("ERROR: Can't read chat combinations from redis for ", c)
		return 
	} else {
		if len(list) == 0 {
			r.rdb.LPush(c, withUserId)
		} else {
			// check if the user exists in the list or not 
			if !utils.Contains(list, withUserId) {
				r.rdb.LPush(c, withUserId)
			}
		}	
	}
}

func (r *RedisService) PushMessageToChatList(message core.DirectMessagePayload) {
	// create combinations based on which userid is smaller 
	key := CreateKeyCombination(message.Sender, message.Receiver)
	if messageJson, err  := json.Marshal(message); err != nil {
		log.Println("ERROR: Unable to marshal the error to push into redis", err)
		return 
	} else {
		if redisErr := r.rdb.LPush(key, messageJson).Err(); redisErr != nil {
			log.Println("ERROR: Unable to marshal the error to push into redis", err)
			return 
		}
	}
}

func CreateKeyCombination(fromUser string, toUser string) string {
	rr, _ := strconv.Atoi(fromUser)
	sr, _ := strconv.Atoi(toUser)	
	if rr > sr {
		return fmt.Sprintf("%v.%v", rr, sr)	
	} else {
		return fmt.Sprintf("%v.%v", sr, rr)	
	}
}
