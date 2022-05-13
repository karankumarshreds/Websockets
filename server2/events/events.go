package events

type EventName string

// Types of events
const (
	NEW_USER       EventName = "NEW_USER"
	DIRECT_MESSAGE EventName = "DIRECT_MESSAGE"
	DISCONNECT     EventName = "DISCONNECT"
	DELETED_USER   EventName = "DELETED_USER"
)