package api

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

type Handler struct {
	db *database.Database
}

// RegisterUser registers a new user
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
