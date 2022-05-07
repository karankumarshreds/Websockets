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
 }
************************************/

type RedisService struct {
	rdb *redis.Client
}

func NewRedisService(rdb *redis.Client) *RedisService {
	// create an empty map for the user data 
	type UserNode struct {
		Username string
		Socket *websocket.Conn
	}
	d := make(map[string]UserNode)
	if usersDataDefault, marshalErr := json.Marshal(d); marshalErr != nil {
		log.Fatal("Unable to marshal initial user data to save in redis", marshalErr)
		return nil 
	} else {
		if err := rdb.Set(USERS_DATA, usersDataDefault, 0).Err(); err != nil {
		log.Fatal("Unable to create an intial user data in redis", err)
		return nil 
	} else {
		result, _ := rdb.Get(USERS_DATA).Result()
		log.Println("Successfully created initial user data in redis", result)
	}
	}
	// TODO create an empty map for the chat data
	return &RedisService{rdb}
}

func (r *RedisService) SetUserRedis(username string) {
	log.Printf("The user with the username %v will be set in redis!", username)
	// get the entire map from redis 
	// create a copy of the entire map 
	// add a key value for the new user 
	// set the updated map in redis 
}

func (r *RedisService) RemoveUserRedis() {

}

func (r *RedisService) SaveMessageRedis() {

}

func (r *RedisService) GetChatRedis() {

}
