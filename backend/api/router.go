package api

import (
	"net/http"
	"real-time-forum/backend/database"
	"time"
)

func NewRouter(db *database.Database, wsHub *Hub) *http.ServeMux {
	r := http.NewServeMux()
	h := Handler{db: db, wsHub: wsHub}
	mw := Middleware{db: db}
	th := NewThrottle(3 * time.Second)

	wrap := func(h http.Handler) http.Handler {
		return mw.LogMiddleware(mw.CorsMiddleware(h))
	}

	r.Handle("/api/check-auth", wrap(mw.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))))

	r.Handle("/api/register", wrap(http.HandlerFunc(h.RegisterUser)))
	r.Handle("/api/login", wrap(http.HandlerFunc(h.LoginUser)))
	r.Handle("/api/logout", wrap(mw.AuthMiddleware(http.HandlerFunc(h.LogoutUser))))

	r.Handle("/api/get-posts", wrap(mw.AuthMiddleware(http.HandlerFunc(h.GetPosts))))
	r.Handle("/api/get-comments", wrap(mw.AuthMiddleware(http.HandlerFunc(h.GetComments))))
	r.Handle("/api/create-post", wrap(th.Throttle(mw.AuthMiddleware(http.HandlerFunc(h.CreatePost)))))
	r.Handle("/api/create-comment", wrap(th.Throttle(mw.AuthMiddleware(http.HandlerFunc(h.CreateComment)))))

	r.Handle("/api/get-users", wrap(mw.AuthMiddleware(http.HandlerFunc(h.GetUsers))))
	r.Handle("/api/messages/{id}", wrap(mw.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.GetMessages(w, r)
		} else if r.Method == "POST" {
			th.Throttle(http.HandlerFunc(h.SendMessage)).ServeHTTP(w, r)
		}
	}))))

	go wsHub.StartHub()
	r.Handle("/api/ws", wrap(mw.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.wsHub.HandleWebSocket(w, r, db)
	}))))

	r.Handle("/", http.FileServer(http.Dir("../frontend")))
	return r
}
