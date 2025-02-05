package utils

import (
	"errors"
	"net/http"

	"app/backend/db"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with its hash.
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateSessionToken creates a new session token using github.com/gofrs/uuid/v5.
func GenerateSessionToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

// GetCookie retrieves the value of the specified cookie.
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// CheckSessionToken validates that the session token from the cookie exists in the online_status table.
func CheckSessionToken(r *http.Request, database *db.Database) error {
	token, err := GetCookie(r, "session_token")
	if err != nil {
		return err
	}
	query := `SELECT id FROM online_status WHERE token = ?`
	row := database.DB.QueryRow(query, token)
	var id int
	if err := row.Scan(&id); err != nil {
		return errors.New("invalid session token")
	}
	return  nil
}
