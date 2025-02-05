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

	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	user.ID = int(userID)
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

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
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		}
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a new session token using github.com/gofrs/uuid/v5.
	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	upsertQuery := `
		INSERT INTO online_status (user_id, token, online, last_seen)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id) DO UPDATE SET
			token = excluded.token,
			online = excluded.online,
			last_seen = excluded.last_seen;
	`
	_, err = h.db.DB.Exec(upsertQuery, user.ID, sessionToken, true)
	if err != nil {
		http.Error(w, "Failed to update online status", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	var post db.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	post.UserID = userID

	query := `INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)`
	result, err := h.db.DB.Exec(query, post.UserID, post.Title, post.Content, post.Category)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}

	post.ID = int(postID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	var posts []db.Post
	query := `SELECT id, user_id, title, content, category FROM posts WHERE user_id = ?`
	rows, err := h.db.DB.Query(query, userID)
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
