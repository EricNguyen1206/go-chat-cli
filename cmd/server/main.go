package main

import (
	"fmt"
	"log"
	"net/http"

	"go-chat-cli/server"
)

func main() {
	hub := server.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})

	addr := ":8080"
	fmt.Println("🚀 The server is running at the port", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
