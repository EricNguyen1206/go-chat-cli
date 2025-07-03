package server

import (
	"fmt"
	"go-chat-cli/utils"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			utils.Info("📥 Client joined:", client.username)
			joinMsg := fmt.Sprintf("🔔 %s has joined the chat room", client.username)
			h.broadcast <- []byte(joinMsg)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				leaveMsg := fmt.Sprintf("🔔 %s has left the chat room", client.username)
				h.broadcast <- []byte(leaveMsg)
				delete(h.clients, client)
				close(client.send)
				utils.Info("❌ Client left:", client.username)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
