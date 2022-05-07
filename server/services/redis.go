package services

import (
	"encoding/json"
	"log"

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

func (r *RedisService) SetUserRedis(userid string, username string, socket *websocket.Conn) {
	log.Printf("Setting user with the username %v in redis!", username)
	userNode := UserNode{
		Username: username,
		Socket: socket,
	}

	// create a map in redis with the name of the userid
	if jsonData, marshalErr  := json.Marshal(userNode); marshalErr != nil {
		log.Println("ERR: Unable to marshal user data to set a map in redis", marshalErr)
		return 
	} else {
		if err := r.rdb.Set(userid, jsonData, 0).Err(); err != nil {
			log.Println("ERR: Unable to set the data in the redis for new user", err)
			return 
		}
	}
	log.Println("Successfully set the data for the user in redis map!")

	// Temporary code below 
	log.Println("Getting the data for the user to see if it was set properly")
	u := r.GetUserRedis(userid)
	log.Println("Got the the user data as", u)
	if err := u.Socket.WriteJSON([]byte("Successfully registerd")); err != nil {
		log.Println("WAIT WHAT??", err)
	}
}

func (r *RedisService) GetUserRedis(userid string) UserNode {
	log.Printf("Getting the user with userId %v", userid)
	var userNode UserNode

	// get the data from the redis database using key of userid 
	if result, err := r.rdb.Get(userid).Result(); err != nil {
		log.Println("ERR: Unable to get the value from redis for userid ", userid, err)
	} else {
		if unmarshalErr := json.Unmarshal([]byte(result), &userNode); unmarshalErr != nil {
			log.Println("ERR: Unable to unmarshal user data from redis db", unmarshalErr)
		}
	}
	return userNode
}

func (r *RedisService) RemoveUserRedis() {

}

func (r *RedisService) SaveMessageRedis() {

}

func (r *RedisService) GetChatRedis() {

}
