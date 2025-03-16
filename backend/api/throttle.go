package api

import (
	"net/http"
	"sync"
	"time"
)

type Throttle struct {
	delay    time.Duration
	lastTime sync.Map
}

func NewThrottle(delay time.Duration) *Throttle {
	return &Throttle{delay: delay}
}

func (t *Throttle) Throttle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		if lastTime, ok := t.lastTime.Load(ip); ok {
			if now.Sub(lastTime.(time.Time)) < t.delay {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}
		t.lastTime.Store(ip, now)

		next.ServeHTTP(w, r)
	})
}
