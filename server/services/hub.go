package service

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}