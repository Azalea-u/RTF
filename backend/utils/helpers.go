package utils

import (
	"database/sql"
	"net/http"
	"real-time-forum/backend/database"
)

// GetCookie retrieves the value of the specified cookie.
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetUserID retrieves the user ID associated with the provided token.
func GetUserID(database *database.Database, token string) (string, error) {
	query := `SELECT id FROM online_status WHERE token = ?`
	row := database.DB.QueryRow(query, token)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil 
		}
		return "", err
	}

	return userID, nil
}
