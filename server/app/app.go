package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"private-chat/handlers"
	"private-chat/services"

	"github.com/go-redis/redis"
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
	// creating a hub for all the users 
	hub := services.NewHub()
	h := handlers.NewHandlers(hub)
	r := mux.NewRouter()

	r.HandleFunc("/", h.HomeHandler).Methods("GET")
	r.HandleFunc("/ws/{userid}/{username}", h.NewWebsocketConnection)

	log.Printf("Server starting at %v", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r)
}


func (a *App) InitRedis(options struct { 
	host string 
	port int 
	password string 
	} ) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(options.host + fmt.Sprint(options.port)),
		DB: 0, // using default db 
		Password: options.password,
	})
	rdb.Set()
	return rdb
}

