package services

const (
	USERS_DATA       string = "USERS_DATA"
	ONLINE_USERS     string = "ONLINE_USERS"
	ONLINE_USERS_MAP string = "ONLINE_USERS_MAP"
)

type redisKeys struct {
	online_users     string
	online_users_map string
}

var REDIS_KEYS redisKeys = redisKeys{
	online_users:     USERS_DATA,
	online_users_map: ONLINE_USERS_MAP,
}
