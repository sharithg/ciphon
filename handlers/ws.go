package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

func (env *Env) Echo(w http.ResponseWriter, r *http.Request) {
	connection, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer connection.Close()

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		err = connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}

		go messageHandler(message)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
