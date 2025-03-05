package api

import (
	"net/http"
	"real-time-forum/backend/database"
)

func NewRouter(db *database.Database , wsHub *Hub) *http.ServeMux {
	r := http.NewServeMux()
	mw := Middleware{db: db}

	wrap := func(h http.Handler) http.Handler {
		return mw.LogMiddleware(mw.CorsMiddleware(mw.AuthMiddleware(h)))
	}

	r.Handle("/", http.FileServer(http.Dir("../frontend")))
	return r
}
