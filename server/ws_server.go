package server

import (
	"net/http"

	"github.com/EricNguyen1206/go-chat-cli/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		utils.Error("Missing username")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Error("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		utils.Info("🔌 Disconnected:", c.username)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			utils.Error("❗ Read error:", err)
			break
		}
		wrapped := []byte("[" + c.username + "]: " + string(message))
		c.hub.broadcast <- wrapped
	}
}


func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c // ✅ trigger remove client
		c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			utils.Error("❗ Write error:", err)
			break
		}
	}
}

