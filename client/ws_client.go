package client

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn *websocket.Conn
	send chan string
	recv chan string
}

func DialWebSocket(url string, timeout time.Duration) (*websocket.Conn, *http.Response, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
	}
	return dialer.Dial(url, nil)
}

func NewWSClient(url string) (*WSClient, error) {
	conn, _, err := DialWebSocket(url, 10*time.Second)
	if err != nil {
		return nil, err
	}

	client := &WSClient{
		conn: conn,
		send: make(chan string),
		recv: make(chan string),
	}

	go client.readPump()
	go client.writePump()

	return client, nil
}

func (c *WSClient) readPump() {
	defer func() {
		c.conn.Close()
		close(c.recv) // Ensure closed once
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("🔌 Read error:", err)
			break
		}

		// ✅ Check if channel is still open
		select {
		case c.recv <- string(message):
		default: // Avoid sending to closed channel
		}
	}
}


func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("❗ Error sending:", err)
			break
		}
	}
}


func (c *WSClient) Close() {
	c.conn.Close()
	close(c.send)
	close(c.recv)
}
