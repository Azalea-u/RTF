package utils

import (
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
	query := `SELECT id FROM user WHERE token = ?`
	row := database.DB.QueryRow(query, token)
	if row.Err() != nil {
		return "", row.Err()
	}
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
