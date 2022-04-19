package handlers

import (
	"log"
	"net/http"
	"private-chat/middlewares"

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
	userid := mux.Vars(r)["userid"]
	// upgrade the http request to websocket connection 
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Cannot upgrade the http request to websocket", err)
		return 
	}
	log.Println(userid, connection)
}

func  
