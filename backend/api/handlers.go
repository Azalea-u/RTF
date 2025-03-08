package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Handler struct {
	db    *database.Database
	wsHub *Hub
}

/* -------------------- Authentication -------------------- */

// RegisterUser  registers a new user
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if err := utils.ValidateUser(user); err != nil {
		http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	user.ID, err = utils.NewUUID()
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO user (id, username, email, password, first_name, last_name, age, gender) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = h.db.DB.Exec(query, user.ID.String(), user.Username, user.Email, hash, user.FirstName, user.LastName, user.Age, user.Gender)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user.Password = ""

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// LoginUser  logs in a user
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	log.Println(credentials)

	var user database.User
	query := `SELECT id, username, email, password, first_name, last_name, age, gender FROM user WHERE email = ? OR username = ?`
	row := h.db.DB.QueryRow(query, credentials.EmailOrUsername, credentials.EmailOrUsername)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age, &user.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"message": "User  not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := utils.CheckPassword(credentials.Password, user.Password); err != nil {
		http.Error(w, `{"message": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.ID.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	query = `UPDATE user SET token = ? WHERE id = ?`
	_, err = h.db.DB.Exec(query, user.ID.String(), user.ID.String())
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user.Password = ""

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// LogoutUser  logs out a user
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now(),
	})

	query := `UPDATE user SET token = NULL WHERE id = ?`
	_, err = h.db.DB.Exec(query, token)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	h.wsHub.logout(userID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

/* -------------------- Websocket -------------------- */

// GetUsers returns a list of users along with their online status
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, "Error getting session token", http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		log.Println("Error getting user ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	query := `
        SELECT
            u.id,
            u.username
        FROM
            user u
        LEFT JOIN (
            SELECT
                CASE
                    WHEN sender_id = ? THEN receiver_id
                    ELSE sender_id
                END AS user_id,
                MAX(created_at) AS latest_message_time
            FROM message
            WHERE sender_id = ? OR receiver_id = ?
            GROUP BY
                CASE
                    WHEN sender_id = ? THEN receiver_id
                    ELSE sender_id
                END
        ) m ON u.id = m.user_id
        WHERE
            u.id <> ?
        ORDER BY
            m.latest_message_time DESC,
            u.username ASC;
    `

	rows, err := h.db.DB.Query(query, userID, userID, userID, userID, userID)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []database.User
	connectedUsers := make(map[string]bool)

	for client := range h.wsHub.clients {
		connectedUsers[client.id] = true
	}

	for rows.Next() {
		var user database.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		user.Online = connectedUsers[user.ID.String()]
		if !user.Online {
			user.Online = false
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetMessages returns the messages between two users 10 at a time
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, "Error getting session token", http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		log.Println("Error getting user ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	otherUserID := r.URL.Path[len("/api/messages/"):]
	if otherUserID == "" {
		http.Error(w, `{"message": "Other user ID is required"}`, http.StatusBadRequest)
		return
	}

	query := `
		SELECT sender_id, receiver_id, content, created_at 
		FROM message 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY created_at DESC 
		LIMIT 10
	`
	rows, err := h.db.DB.Query(query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		log.Println("Error querying messages:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []database.Message
	for rows.Next() {
		var message database.Message
		if err := rows.Scan(&message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt); err != nil {
			log.Println("Error scanning message:", err)
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	if messages == nil {
		messages = []database.Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println("Error encoding messages:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	}
}

// SendMessage sends a message to another user
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, "Error getting session token", http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		log.Println("Error getting user ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	otherUserID := r.URL.Path[len("/api/messages/"):]
	if otherUserID == "" {
		http.Error(w, `{"message": "Other user ID is required"}`, http.StatusBadRequest)
		return
	}

	var message database.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Println("Error decoding message:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if message.Content == "" || strings.TrimSpace(message.Content) == "" {
		http.Error(w, `{"message": "Message cannot be empty"}`, http.StatusBadRequest)
		return
	}

	message.SenderID, err = uuid.FromString(userID)
	if err != nil {
		log.Println("Error parsing sender ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	message.ReceiverID, err = uuid.FromString(otherUserID)
	if err != nil {
		log.Println("Error parsing receiver ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO message (sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err = h.db.DB.Exec(query, message.SenderID, message.ReceiverID, message.Content, time.Now())
	if err != nil {
		log.Println("Error inserting message:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
