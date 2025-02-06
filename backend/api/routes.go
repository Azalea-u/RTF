package api

import (
	"app/backend/db"
	"net/http"
)

func NewRouter(db *db.Database) *http.ServeMux {
	r := http.NewServeMux()
	h := &Handler{db: db}
	m := &Middleware{db: db}

	// Middleware wrapper for cleaner code
	wrap := func(handler http.Handler) http.Handler {
		return m.LoggingMiddleware(m.CORSMiddleware(handler))
	}

	// API endpoints
	r.Handle("/api/register", wrap(http.HandlerFunc(h.RegisterUser)))
	r.Handle("/api/login", wrap(http.HandlerFunc(h.LoginUser)))
	r.Handle("/api/logout", wrap(m.AuthMiddleware(http.HandlerFunc(h.LogoutUser))))
	r.Handle("/api/create-post", wrap(m.AuthMiddleware(http.HandlerFunc(h.CreatePost))))
	r.Handle("/api/posts", wrap(m.AuthMiddleware(http.HandlerFunc(h.GetPosts))))
	r.Handle("/api/user", wrap(m.AuthMiddleware(http.HandlerFunc(h.GetUserData))))

	// Messaging Endpoints
	r.Handle("/api/users", wrap(m.AuthMiddleware(http.HandlerFunc(h.GetAllUsers))))         // Fetch user list
	r.Handle("/api/messages", wrap(m.AuthMiddleware(http.HandlerFunc(h.GetMessages))))      // Get messages
	r.Handle("/api/messages", wrap(m.AuthMiddleware(http.HandlerFunc(h.SendMessage))))      // Send message

	// WebSocket Chat Endpoint
	wsHub := NewHub()
	go wsHub.Run()
	r.Handle("/api/chat", wrap(m.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsHub, w, r)
	}))))

	// Static files
	r.Handle("/", wrap(http.FileServer(http.Dir("../frontend"))))

	return r
}
