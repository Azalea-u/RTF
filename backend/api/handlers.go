package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
	"strconv"
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

	token, err := utils.NewUUID()
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	query = `UPDATE user SET token = ? WHERE id = ?`
	_, err = h.db.DB.Exec(query, token, user.ID)
	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user.Password = ""

	h.wsHub.login(user.ID.String())

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

	if users == nil {
		users = []database.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetMessages returns all messages between two users
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

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	query := `
		SELECT sender_id, receiver_id, content, created_at 
		FROM message 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := h.db.DB.Query(query, userID, otherUserID, otherUserID, userID, limit, offset)
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

	// broadcast message to other user
	h.wsHub.broadcast <- []byte(`{"type": "message", "content": "` + message.SenderID.String() + `,` + message.ReceiverID.String() + `"}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/* -------------------- Posts&Comments -------------------- */

// GetPosts gets all posts
func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []database.Post

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	query := `
        SELECT id, user_id, title, content, category, created_at
        FROM post
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := h.db.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("Failed to fetch posts: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post database.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil {
			log.Printf("Failed to scan post: %v", err)
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate over posts: %v", err)
		http.Error(w, "Failed to iterate over posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(posts) == 0 {
		posts = []database.Post{}
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println("Error encoding posts:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	}
}

// CreatePost creates a new post
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
    userID, err := utils.GetUserID(h.db, token)
    if err != nil {
        log.Println("Error getting user ID:", err)
        http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

    var post database.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        log.Println("Error decoding post:", err)
        http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

    err = utils.ValidatePost(post)
    if err != nil {
        log.Println("Error validating post:", err)
        http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

    post.UserID, err = uuid.FromString(userID)
    if err != nil {
        log.Println("Error parsing user ID:", err)
        http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

	log.Println("Post:", post)

    query := `INSERT INTO post (user_id, title, content, category, created_at) VALUES (?, ?, ?, ?, ?)`
    _, err = h.db.DB.Exec(query, post.UserID, post.Title, post.Content, post.Category, time.Now())
    if err != nil {
        log.Println("Error inserting post:", err)
        http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

// GetComments gets all comments for a post
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, `{"message": "Invalid post ID"}`, http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	var comments []database.Comment

	query := `
        SELECT id, content, user_id, post_id, created_at
        FROM comment
        WHERE post_id = ?
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := h.db.DB.Query(query, postID, limit, offset)
	if err != nil {
		log.Printf("Failed to fetch comments: %v", err)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment database.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt); err != nil {
			log.Printf("Failed to scan comment: %v", err)
			http.Error(w, "Failed to scan comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate over comments: %v", err)
		http.Error(w, "Failed to iterate over comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(comments) == 0 {
		comments = []database.Comment{}
	}
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Println("Error encoding comments:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	}
}

// CreateComment creates a new comment
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		log.Println("Error getting session token:", err)
		http.Error(w, "Error getting session token", http.StatusInternalServerError)
		return
	}
	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		log.Println("Error getting user ID:", err)
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	var comment database.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Println("Error decoding comment:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	
	if comment.Content == "" || strings.TrimSpace(comment.Content) == "" {
		http.Error(w, `{"message": "Comment content is required"}`, http.StatusBadRequest)
		return
	}

	comment.UserID, err = uuid.FromString(userID)
	if err != nil {
		log.Println("Error parsing user ID:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO comment (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err = h.db.DB.Exec(query, comment.PostID, comment.UserID, comment.Content, time.Now())
	if err != nil {
		log.Println("Error inserting comment:", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}