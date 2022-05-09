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

type RedisConfig struct {
	host string 
	port int 
	password string 
}

func (a *App) InitRoutes() {
	// creating a hub for all the users 
	rdb := a.InitRedis(RedisConfig{
		host: "localhost",
		port: 6379, 
		password: "",
	})
	
	hubService    := services.NewHub(rdb)
	redisService  := services.NewRedisService(rdb)

	h := handlers.NewHandlers(hubService, redisService)
	r := mux.NewRouter()


	r.HandleFunc("/", h.HomeHandler).Methods("GET")
	r.HandleFunc("/ws/{userid}/{username}", h.NewWebsocketConnection)
	r.HandleFunc("/get-chats/{userid}", h.GetAllChats).Methods("GET")

	log.Printf("Server starting at %v", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r)
}


func (a *App) InitRedis(options RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(options.host + ":" + fmt.Sprint(options.port)),
		DB: 0, // using default db 
		Password: options.password,
	})
	return rdb
}

