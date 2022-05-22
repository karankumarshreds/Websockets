package services

const (
	USERS_DATA   string = "USERS_DATA"
	ONLINE_USERS string = "ONLINE_USERS"
)

var REDIS_KEYS struct{ online_users string } = struct{ online_users string }{
	online_users: USERS_DATA,
}
