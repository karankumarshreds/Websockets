package core

import "private-chat/events"

// Standard structure of how the payload will look like
type EventPayload struct {
	EventName    events.EventName `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

/************************* NOTE ***************************/
/* Each event payload must have a unique userId property  */
/**********************************************************/

// Payload structure for new user event
type NewUserPayload struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// Payload structure for direct message event
type DirectMessagePayload struct {
	UserId   string `json:"userId"`
	Sender   string `json:"sender"`
	Message  string `json:"message"`
}

// Payload structure for disconnect event 
type DisconnectPayload struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// "Response" Payload structure for direct message
type DirectMessageResponse struct {
	UserId string `json:"userId"`
	Message string `json:"message"`
	Sender string `json:"sender"`
	Time string `json:"time"`
}

// REDIS RELATED 

// key denotes the userid 
type OnlineUsersRedisMap map[string]NewUserPayload 
