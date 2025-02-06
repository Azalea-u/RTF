package api

import (
	"app/backend/db"
	"app/backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Handler struct {
	db *db.Database
}

// Helper function: Write JSON Response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function: Retrieve User ID from Session
func (h *Handler) getUserIDFromSession(r *http.Request) (int, error) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		return 0, err
	}

	var userID int
	query := `SELECT user_id FROM online_status WHERE token = ? AND online = TRUE`
	row := h.db.DB.QueryRow(query, token)
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("invalid session token")
		}
		return 0, err
	}
	return userID, nil
}

// User Registration
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (username, email, password, first_name, last_name, gender) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := h.db.DB.Exec(query, user.Username, user.Email, hashedPassword, user.FirstName, user.LastName, user.Gender)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()
	user.ID = int(userID)
	user.Password = ""
	writeJSON(w, http.StatusCreated, user)
}

// User Login
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user db.User
	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := h.db.DB.QueryRow(query, credentials.Username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionToken, _ := utils.GenerateSessionToken()
	h.db.DB.Exec(`INSERT INTO online_status (user_id, token, online, last_seen)
                  VALUES (?, ?, TRUE, CURRENT_TIMESTAMP)
                  ON CONFLICT(user_id) DO UPDATE SET token = excluded.token, online = TRUE, last_seen = CURRENT_TIMESTAMP;`,
		user.ID, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}

// User Logout
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetCookie(r, "session_token")
	h.db.DB.Exec(`UPDATE online_status SET online = FALSE WHERE token = ?`, token)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

// Create Post
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post db.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	post.UserID = userID
	result, _ := h.db.DB.Exec(`INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)`, post.UserID, post.Title, post.Content, post.Category)
	postID, _ := result.LastInsertId()
	post.ID = int(postID)

	writeJSON(w, http.StatusCreated, post)
}
// Get Posts
func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []db.Post
	query := `SELECT id, user_id, title, content, category FROM posts ORDER BY created_at DESC`
	rows, err := h.db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to iterate over posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// Get User List (For Messaging)
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.db.DB.Query(`SELECT id, username FROM users`)
	defer rows.Close()

	var users []db.User
	for rows.Next() {
		var user db.User
		rows.Scan(&user.ID, &user.Username)
		users = append(users, user)
	}

	writeJSON(w, http.StatusOK, users)
}
// GetMessages retrieves messages between two users.
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ReceiverID int `json:"receiver_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	senderID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch conversation ID
	var conversationID int
	err = h.db.DB.QueryRow(`SELECT id FROM conversations 
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)`, 
		senderID, request.ReceiverID, request.ReceiverID, senderID).Scan(&conversationID)

	if err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Fetch messages in conversation
	rows, err := h.db.DB.Query(`SELECT id, conversation_id, sender_id, content, read, created_at FROM messages 
		WHERE conversation_id = ? ORDER BY created_at ASC`, conversationID)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []db.Message
	for rows.Next() {
		var msg db.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.Read, &msg.CreatedAt); err != nil {
			http.Error(w, "Failed to parse messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	writeJSON(w, http.StatusOK, messages)
}

// SendMessage allows a user to send a message to another user.
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ReceiverID int    `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	senderID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure a conversation exists
	var conversationID int
	err = h.db.DB.QueryRow(`SELECT id FROM conversations 
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)`, 
		senderID, request.ReceiverID, request.ReceiverID, senderID).Scan(&conversationID)

	if err == sql.ErrNoRows {
		// Create a new conversation if none exists
		result, err := h.db.DB.Exec(`INSERT INTO conversations (user1_id, user2_id, created_at) VALUES (?, ?, ?)`,
			senderID, request.ReceiverID, time.Now())
		if err != nil {
			http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
			return
		}
		conversationID64, _ := result.LastInsertId()
		conversationID = int(conversationID64)
	} else if err != nil {
		http.Error(w, "Failed to retrieve conversation", http.StatusInternalServerError)
		return
	}

	// Insert the new message
	_, err = h.db.DB.Exec(`INSERT INTO messages (conversation_id, sender_id, content, read, created_at) VALUES (?, ?, ?, ?, ?)`,
		conversationID, senderID, request.Content, false, time.Now())

	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Message sent"})
}

// Get User Data
func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	var user db.User
	query := `SELECT id, username, email, first_name, last_name, gender FROM users WHERE id = ?`
	row := h.db.DB.QueryRow(query, userID)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Gender); err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// Exclude password for security
	user.Password = ""

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}