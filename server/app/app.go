package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type App struct {}

func NewApp () *App {
	return &App{}
}

func (a *App) Run () {
	godotenv.Load()
	fmt.Printf("Server starting at %v", os.Getenv("PORT"))
}