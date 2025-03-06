package api

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
	"time"
)

type Middleware struct {
	db *database.Database
}

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("\033[36m--------------------------------------\033[0m")
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Completed request: %s %s in %v", r.Method, r.URL.Path, duration)
	})
}

func (m *Middleware) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Specify allowed headers
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.GetCookie(r, "session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := utils.GetUserID(m.db, token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", userID)
		next.ServeHTTP(w, r)
	})
}
