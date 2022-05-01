package services

import "log"

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    map[*Client]bool{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run acts like an interface between the readPump and writePump and updates Hub map
func (h *Hub) Run() {
	for { // infinite loop
		select {
		case client := <-h.Register:
			log.Println("Registering user with userid", client.UserId)
			h.Clients[client] = true
			onlineUsers := []string{}
			for c := range h.Clients {
				onlineUsers = append(onlineUsers, c.UserId)
			}
			log.Println("List of online users", onlineUsers)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		}
	}
}