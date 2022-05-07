package services

import "github.com/go-redis/redis"

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
	// create an empty map for the chat data 
	return &RedisService{rdb}
}

func (r *RedisService) SetUserRedis() {
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
