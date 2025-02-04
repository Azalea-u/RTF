package api

import (
	"app/backend/db"
	"app/backend/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

// Middleware struct holds dependencies for middleware functions
type Middleware struct {
	db *db.Database
}

// AuthMiddleware ensures the request is authenticated
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Printf("ERROR: Missing session token - %v", err)
			http.Error(w, "Unauthorized: Missing session token", http.StatusUnauthorized)
			return
		}

		// Validate the session token against the database
		var userID int
		query := `SELECT user_id FROM online_status WHERE token = ? AND online = TRUE`
		row := m.db.DB.QueryRow(query, cookie.Value)
		if err := row.Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("ERROR: Invalid session token - %v", err)
				http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
			} else {
				log.Printf("ERROR: Failed to validate session - %v", err)
				http.Error(w, "Failed to validate session", http.StatusInternalServerError)
			}
			return
		}

		// Attach the user ID to the request context
		ctx := utils.SetUserIDInContext(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs incoming requests
func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware enables CORS for all routes
func (m *Middleware) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}