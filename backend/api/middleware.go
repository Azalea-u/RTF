package api

import (
	"app/backend/db"
	"app/backend/utils"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	db *db.Database
}

// LoggingMiddleware logs the request method and URL.
func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v \n", r.URL.Path, time.Since(start))
	})
}

// CORSMiddleware adds CORS headers to the response.
func (m *Middleware) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware verifies the JWT token in the Authorization header.
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.GetCookie(r, "session_token")
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		if err := utils.CheckSessionToken(r, m.db); err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
