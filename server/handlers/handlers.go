package handlers

import (
	"log"
	"net/http"
	"private-chat/core"
	"private-chat/middlewares"
	"private-chat/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Handlers struct {}

var upgrader websocket.Upgrader = websocket.Upgrader{
		ReadBufferSize	: 	1024,
		WriteBufferSize	: 	1024,
		CheckOrigin: func(r *http.Request) bool {return true},
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	middlewares.JsonResponse(w,r,http.StatusOK, "Welcome")
}

func (h *Handlers) NewWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	hub := services.NewHub()
	go hub.Run()
	userid := mux.Vars(r)["userid"]
	username := mux.Vars(r)["username"]
	// upgrade the http request to websocket connection 
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Cannot upgrade the http request to websocket", err)
		return 
	}
	log.Println("Creating a new websocket connection for user", userid)
	// Creating a new user struct
	client := &services.Client{
		UserId: userid,
		Username: username,
		Hub: hub,
		Conn: connection,
		Send: make(chan core.EventPayload),
	}

	// Registering the user to the hub
	client.Hub.Register <- client 

	go client.ReadPump()
	go client.WritePump()
	
}


