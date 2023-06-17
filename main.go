package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	clients := make(map[*websocket.Conn]bool)
	broadcast := make(chan Message)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		clients[conn] = true

		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println(err)
				delete(clients, conn)
				break
			}

			broadcast <- message
		}
	})

	go func() {
		ticker := time.NewTicker(3 * time.Second) // Enviar mensaje cada 3 segundos

		for range ticker.C {
			message := Message{Message: "Información recurrente"}

			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
