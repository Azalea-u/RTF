package api

import (
	"app/backend/db"
	"net/http"
)

func NewRouter(db *db.Database) *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}
