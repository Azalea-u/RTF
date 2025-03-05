package api

import (
	"net/http"
	"real-time-forum/backend/database"
)

func NewRouter(db *database.Database , wsHub *Hub) *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/", http.FileServer(http.Dir("../frontend")))
	return r
}
