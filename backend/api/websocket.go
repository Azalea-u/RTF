package api

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	id   string
}

type Hub struct {
	clients    map[*Client]bool // Registered clients connected to the hub.
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub initializes a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// StartHub starts the hub to manage WebSocket connections.
func (h *Hub) StartHub() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client registered:", client.id)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				_ = client.conn.Close()
				log.Println("Client unregistered:", client.id)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("Error sending message to client:", client.id, err)
					_ = client.conn.Close()
					h.unregister <- client
				}
			}
		}
	}
}

// HandleWebSocket handles WebSocket connections.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, db *database.Database) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		http.Error(w, "Error upgrading to WebSocket", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, "Error getting session token", http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(db, token)
	if err != nil {
		log.Println("Error getting user ID:", err)
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	client := &Client{conn: ws, id: userID}
	h.register <- client
	defer func() {
		h.unregister <- client
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		h.broadcast <- message
	}
}

func (h *Hub) Shutdown() {
	for client := range h.clients {
		_ = client.conn.Close()
		h.unregister <- client
	}
	log.Println("WebSocket hub shutdown completed.")
}

func (h *Hub) logout(userID string) {
	for client := range h.clients {
		if client.id == userID {
			_ = client.conn.Close()
			h.unregister <- client
		}
	}
}