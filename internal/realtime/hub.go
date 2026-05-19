package realtime

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)


type Hub struct {
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)

	if(err != nil) {
		return 
	}

	h.clients[con] = true
}

func (h *Hub) Broadcast(event any) {
	data, err := json.Marshal(event)

	if(err != nil) {
		return 
	}

	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			client.Close()
			delete(h.clients, client)
		}
	}
}