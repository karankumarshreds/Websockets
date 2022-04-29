package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"private-chat/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {}

func NewApp () *App {
	return &App{}
}

func (a *App) Run () {
	godotenv.Load()
	a.InitRoutes()
}

func (a *App) InitRoutes() {
	h := handlers.NewHandlers()
	r := mux.NewRouter()

	r.HandleFunc("/", h.HomeHandler).Methods("GET")
	r.HandleFunc("/ws/{userid}", h.NewWebsocketConnection)

	log.Printf("Server starting at %v", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r)
}

