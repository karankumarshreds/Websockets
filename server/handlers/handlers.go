package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"private-chat/core"
	"private-chat/middlewares"
	"private-chat/services"

	"github.com/go-redis/redis"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	hub *services.Hub
	redisService *services.RedisService
	rdb *redis.Client
}

var upgrader websocket.Upgrader = websocket.Upgrader{
		ReadBufferSize	: 	1024,
		WriteBufferSize	: 	1024,
		CheckOrigin: func(r *http.Request) bool {return true},
}

func NewHandlers(hub *services.Hub, redisService *services.RedisService, rdb *redis.Client) *Handlers {
	return &Handlers{hub, redisService, rdb}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	middlewares.JsonResponse(w,r,http.StatusOK, "Welcome")
}

func (h *Handlers) NewWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	go h.hub.Run()
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
	client := services.NewClientService(
		h.hub, 
		connection, 
		make(chan core.EventPayload), 
		userid, 
		username, 
		h.redisService,
	)

	// Registering the user to the hub
	client.Hub.Register <- client 

	go client.ReadPump()
	go client.WritePump()

	/* start listening for external messages */
	l := services.NewListeners(h.rdb, h.hub)
	go l.NewUserListener()
	go l.DirectMessageListener()
	
}

func (h *Handlers) GetAllChats(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["userid"]
	response := h.redisService.GetAllChatsWithLastMessage(userid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

