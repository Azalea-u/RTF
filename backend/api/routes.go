package api

import (
	"app/backend/db"
	"net/http"
)

func NewRouter(db *db.Database) *http.ServeMux {
	r := http.NewServeMux()
	h := &Handler{db: db}
	m := &Middleware{db: db}

	// API endpoints with logging, CORS, and authentication where needed.
	r.Handle("/api/register", m.LoggingMiddleware(m.CORSMiddleware(http.HandlerFunc(h.RegisterUser))))
	r.Handle("/api/login", m.LoggingMiddleware(m.CORSMiddleware(http.HandlerFunc(h.LoginUser))))
	r.Handle("/api/logout", m.LoggingMiddleware(m.CORSMiddleware(m.AuthMiddleware(http.HandlerFunc(h.LogoutUser)))))
	r.Handle("/api/create-post", m.LoggingMiddleware(m.CORSMiddleware(m.AuthMiddleware(http.HandlerFunc(h.CreatePost)))))
	r.Handle("/api/posts", m.LoggingMiddleware(m.CORSMiddleware(m.AuthMiddleware(http.HandlerFunc(h.GetPosts)))))
	r.Handle("/api/user", m.LoggingMiddleware(m.CORSMiddleware(m.AuthMiddleware(http.HandlerFunc(h.GetUserData)))))

	// Static files
	fs := http.FileServer(http.Dir("../frontend"))
	r.Handle("/", m.LoggingMiddleware(m.CORSMiddleware(fs)))

	return r
}
