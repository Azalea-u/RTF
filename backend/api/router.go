package api

import (
	"net/http"
	"real-time-forum/backend/database"
)

func NewRouter(db *database.Database , wsHub *Hub) *http.ServeMux {
	r := http.NewServeMux()
	h := Handler{db: db, wsHub: wsHub}
	mw := Middleware{db: db}

	wrap := func(h http.Handler) http.Handler {
		return mw.LogMiddleware(mw.CorsMiddleware(h))
	}

	r.Handle("/api/check-auth", wrap(mw.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))))

	r.Handle("/api/register", wrap(http.HandlerFunc(h.RegisterUser)))
	r.Handle("/api/login", wrap(http.HandlerFunc(h.LoginUser)))
	r.Handle("/api/logout", wrap(mw.AuthMiddleware(http.HandlerFunc(h.LogoutUser))))

	r.Handle("/api/get-users", wrap(mw.AuthMiddleware(http.HandlerFunc(h.GetUsers))))
	r.Handle("/api/messages/{id}", wrap(mw.AuthMiddleware(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.GetMessages(w, r)
		} else if r.Method == "POST" {
			h.SendMessage(w, r)
		}
	}))))

	r.Handle("/", http.FileServer(http.Dir("../frontend")))
	return r
}
