package handlers

import "github.com/gorilla/mux"

type Handlers struct {}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) InitRoutes() {
	r := mux.NewRouter()
}


