package api

import (
	"app/backend/db"
	"net/http"
)

func NewRouter(db *db.Database) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		(&Handler{db}).RegisterUser(w, r)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		(&Handler{db}).LoginUser(w, r)
	})
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		(&Handler{db}).CreatePost(w, r)
	})
	return mux
}
