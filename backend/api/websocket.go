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
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) StartHub() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client registered:", client.id)
			h.broadcast <- []byte(`{"type": "user_connected"}`)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				_ = client.conn.Close()
				log.Println("Client unregistered:", client.id)
				h.broadcast <- []byte(`{"type": "user_disconnected"}`)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					_ = client.conn.Close()
					log.Println("Client forcefully disconnected due to blocked send channel:", client.id)
				}
			}
		}
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, db *database.Database) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		http.Error(w, "Error upgrading to WebSocket", http.StatusInternalServerError)
		return
	}

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

	client := &Client{
		conn: ws,
		id:   userID,
		send: make(chan []byte, 256),
	}

	h.register <- client
	go client.writePump()
	defer func() {
		h.unregister <- client
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		h.broadcast <- []byte(`{"type": "message", "content": "` + string(message) + `"}`)
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("WebSocket write error for client:", c.id, err)
			break
		}
	}
	c.conn.Close()
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
			h.broadcast <- []byte(`{"type": "user_disconnected"}`)
		}
	}
}

func (h *Hub) login(userID string) {
	for client := range h.clients {
		if client.id == userID {
			h.broadcast <- []byte(`{"type": "user_connected"}`)
		}
	}
}
