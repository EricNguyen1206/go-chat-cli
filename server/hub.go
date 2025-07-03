package server

import (
	"fmt"
	"github.com/EricNguyen1206/go-chat-cli/utils"
)

type Hub struct {
	clients    map[string]*Client // map username → client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if oldClient, exists := h.clients[client.username]; exists {
				utils.Warn("⚠️ Duplicate user detected, disconnecting old client:", client.username)
				oldClient.Disconnect() // 🔄 Trigger unregister cleanly
			}
			h.clients[client.username] = client
			utils.Info("📥 Client joined:", client.username)
		

		case client := <-h.unregister:
			if existing, ok := h.clients[client.username]; ok && existing == client {
				leaveMsg := fmt.Sprintf("🔔 %s has left the chat room", client.username)
				h.broadcast <- []byte(leaveMsg)
				delete(h.clients, client.username)
				close(client.send)
				utils.Info("❌ Client left:", client.username)
			}	

		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					utils.Warn("⚠️ Failed to send to", client.username)
					close(client.send)
					delete(h.clients, client.username)
				}
			}
		}
	}
}
