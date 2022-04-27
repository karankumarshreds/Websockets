package core

import "private-chat/events"

// Standard structure of how the payload will look like
type EventPayload struct {
	EventName    events.EventName `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

// Payload structure for new user event
type NewUserPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

// Payload structure for direct message event
type DirectMessagePayload struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

// Payload structure for disconnect event 
type DisconnectPayload struct {
	Username string `json:"username"`
	UserId string `json:"userId"`
}
