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
	log.Println("Initiating the chat save process on redis")
	r.CreateChatCombinations(message)
	r.PushMessageToChatList(message)
}

func (r *RedisService) GetAllChatsWithLastMessage(receiver string) *[]core.DirectMessagePayload {
	// get the people the user has chatted with using the key <receiver>_CHATS 
	chattedWith, err := r.rdb.LRange(fmt.Sprintf("%v_CHATS", receiver), 0, -1).Result()
	if err != nil {
		log.Println("ERROR: Cannot get list of chatted with for user", receiver)
		return nil
	}

	// response chats with last message 
	var chats []core.DirectMessagePayload

	// check for all the combinations using the same logic using CreateKeyCombination function 
	for _, user := range chattedWith {
		key := CreateKeyCombination(user, receiver)
		chat, _ := r.rdb.LRange(key, 0, 0).Result()
		if len(chat) == 1 {
			var c core.DirectMessagePayload
			if err := json.Unmarshal([]byte(chat[0]), &c); err != nil {
				log.Println("ERROR: Unable to unmarshal chat for users combination of ", key)
				return nil
			} else {
				chats = append(chats, c)	
			}
		}
	}
	
	// return all the arrays of messages using LRANGE using 0 0 to only include the latest message 
	log.Println("Got the list of the chats with last messages", chats)
	return &chats
}

 
func (r *RedisService) GetChatRedis(receiver string) *[]core.DirectMessagePayload {
	var	chats []core.DirectMessagePayload
	// get receiver chats array using <receiver>_CHATS 
	key := fmt.Sprintf("%v_CHATS", receiver)
	if chattedWith, err := r.rdb.LRange(key, 0, -1).Result(); err != nil {
		log.Println("ERROR: Unable to get list of chatted with to create combinations", err)
		return nil 
	} else {
		// create combinations using this chatted with list 
		for _, userid := range chattedWith {
			combination := CreateKeyCombination(userid, receiver)
			if msgs, redisErr := r.rdb.LRange(combination, 0, 0).Result(); redisErr != nil {
				log.Println("ERROR: Unable tot get messages from redis for the user")
				return nil 
			} else {
				lastMessage := msgs[0]
				var msg core.DirectMessagePayload
				if unmarshalErr := json.Unmarshal([]byte(lastMessage), &msg); unmarshalErr != nil {
					log.Println("Unabele to marshal single error ", unmarshalErr)
					return nil
				} else {
					chats = append(chats, msg)
				}
			}
		}
		log.Println("Created chat array with the lates messages as ", chats)
		
	}
	return &chats
}

func (r *RedisService) CreateChatCombinations(message core.DirectMessagePayload)  {
	log.Println("CreateChatCombinations(): Creating chat combination for rr and sr")
	c := fmt.Sprintf("%v_CHATS", message.Receiver)
	if list, err  := r.rdb.LRange(c, 0, -1).Result(); err != nil {
		log.Println("ERROR: Can't read chat combinations from redis for ", c)
		return 
	} else {
		if len(list) == 0 {
			if redisErr := r.rdb.LPush(c, message.Sender).Err(); redisErr != nil {
				log.Panic("ERROR: A: while pushing combination array", redisErr)
			}
		} else {
			// check if the user exists in the list or not 
			if !utils.Contains(list, message.Sender) {
				if redisErr := r.rdb.LPush(c, message.Sender).Err(); redisErr != nil {
					log.Panic("ERROR: B: while pushing combination array", redisErr)
				}
			}
		}	
		log.Println("Created chat combinations as ", c)
	}
}


func (r *RedisService) PushMessageToChatList(message core.DirectMessagePayload) {
	// create combinations based on which userid is smaller 
	log.Println("PushMessageToChatList(): Pushing the message to the combination key")
	key := CreateKeyCombination(message.Sender, message.Receiver)
	log.Println("PushMessageToChatList(): The ss+rr combination key is", key)
	if messageJson, err  := json.Marshal(message); err != nil {
		log.Println("ERROR: Unable to marshal the error to push into redis", err)
		return 
	} else {
		if redisErr := r.rdb.LPush(key, messageJson).Err(); redisErr != nil {
			log.Println("ERROR: Unable to marshal the error to push into redis", err)
			return 
		} 
		log.Println("PushMessageToChatList(): Message successfully saved")
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
