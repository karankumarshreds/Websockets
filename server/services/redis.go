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


type RedisService struct {
	rdb *redis.Client
}

type UserNode struct {
		Username string
		Socket *websocket.Conn
}

type SaveMessageArg struct {
	Sender struct{ Username string; UserId string } 
	Receiver struct { Username string; UserId string }
	Message string 
	Time string
} 

func NewRedisService(rdb *redis.Client) *RedisService {
	return &RedisService{rdb}
}

func (r *RedisService) SaveMessageRedisToChat(message SaveMessageArg) {
	log.Println("Initiating the chat save process on redis")
	r.CreateChatCombinations(message.Sender.UserId, message.Receiver.UserId)
	r.PushMessageToChatList(message)
}

func (r *RedisService) GetAllChatsWithLastMessage(receiver string) *[]SaveMessageArg {
	// get the people the user has chatted with using the key <receiver>_CHATS 
	chattedWith, err := r.rdb.LRange(fmt.Sprintf("%v_CHATS", receiver), 0, -1).Result()
	if err != nil {
		log.Println("ERROR: Cannot get list of chatted with for user", receiver)
		return nil
	}

	// response chats with last message 
	var chats []SaveMessageArg

	// check for all the combinations using the same logic using CreateKeyCombination function 
	for _, user := range chattedWith {
		key := CreateKeyCombination(user, receiver)
		chat, _ := r.rdb.LRange(key, 0, 0).Result()
		if len(chat) == 1 {
			var _c json.RawMessage 
			var c SaveMessageArg

			if err := json.Unmarshal([]byte(chat[0]), &_c); err != nil {
				log.Println("ERROR: Unable to unmarshal to json raw format", err)
				return nil
			} else {
				if err := json.Unmarshal([]byte(_c), &c); err != nil {
					log.Println("ERROR: Unable to unmarshal string to struct", err)
				} else {
					log.Println("Appending a message for receiver", c.Message)
					chats = append(chats, c)	
				}
			}
		}
	}
	// return all the arrays of messages using LRANGE using 0 0 to only include the latest message 
	log.Println("Created the final list of the chats with last messages", chats)
	return &chats
}

func (r *RedisService) CreateChatCombinations(sender string, receiver string)  {
	log.Println("CreateChatCombinations(): Creating chat combination for rr and sr")
	c := fmt.Sprintf("%v_CHATS", receiver)
	if list, err  := r.rdb.LRange(c, 0, -1).Result(); err != nil {
		log.Println("ERROR: Can't read chat combinations from redis for ", c)
		return 
	} else {
		if len(list) == 0 {
			if redisErr := r.rdb.LPush(c, sender).Err(); redisErr != nil {
				log.Panic("ERROR: A: while pushing combination array", redisErr)
			}
		} else {
			// check if the user exists in the list or not 
			if !utils.Contains(list, sender) {
				if redisErr := r.rdb.LPush(c, sender).Err(); redisErr != nil {
					log.Panic("ERROR: B: while pushing combination array", redisErr)
				}
			}
		}	
		log.Println("Created chat combinations as ", c)
	}
}


func (r *RedisService) PushMessageToChatList(msg SaveMessageArg) {
	// create combinations based on which userid is smaller 
	log.Println("PushMessageToChatList(): Pushing the message to the combination key")
	key := CreateKeyCombination(msg.Sender.UserId, msg.Receiver.UserId)
	log.Println("PushMessageToChatList(): The ss+rr combination key is", key)
	if messageJson, err  := json.Marshal(msg); err != nil {
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

func (r *RedisService) RemoveUserFromOnlineMap(userId string) error {
	log.Println("Cleaning up user from the redis map")
	onlineUsersRedisMap, err := r.GetOnlineUsersRedisMap()
	if err != nil {
		return err
	}
	// delete(onlineUsersRedisMap, userId)
	log.Println("Deleted the user from the map, updating redis map again")
	marshalledMap, err := json.Marshal(onlineUsersRedisMap)
	if err != nil {
		log.Println("Unable to marshal the updated map of users", err)
		return err
	} else {
		if err := r.rdb.Set(REDIS_KEYS.online_users_map,marshalledMap,0).Err(); err != nil {
			log.Println("unable to set updated map on redis after deleting user", err)
			return err 
		} 
	}
	return nil
}

func (r *RedisService) GetOnlineUsersRedisMap() (core.OnlineUsersRedisMap, error) {
	var onlineUsersRedisMap core.OnlineUsersRedisMap
	data, err := r.rdb.Get(REDIS_KEYS.online_users_map).Result()
	if err != nil {
		log.Println("Unable to fetch online users map from redis", err)
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &onlineUsersRedisMap); err != nil {
		log.Println("Unable to unmarshal online users map from redis", err)
		return nil, err
	}
	return onlineUsersRedisMap, nil
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


