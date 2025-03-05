package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
	"time"
)

type Handler struct {
	db    *database.Database
	wsHub *Hub
}

// RegisterUser  registers a new user
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.ID, err = utils.NewUUID()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO user (id, username, email, password, first_name, last_name, age, gender) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = h.db.DB.Exec(query, user.ID.String(), user.Username, user.Email, hash, user.FirstName, user.LastName, user.Age, user.Gender)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser  logs in a user
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user database.User
	query := `SELECT id, username, email, password, first_name, last_name, age, gender FROM user WHERE email = ? OR username = ?`
	row := h.db.DB.QueryRow(query, credentials.EmailOrUsername, credentials.EmailOrUsername)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age, &user.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User  not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := utils.CheckPassword(credentials.Password, user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// LogoutUser  logs out a user
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, err := utils.GetUserID(h.db, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.wsHub.logout(userID)
	w.WriteHeader(http.StatusOK)
}
