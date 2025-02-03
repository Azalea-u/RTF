package utils

import (
	"app/backend/db"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func CheckSessionToken(r *http.Request, db *db.Database) (bool, error) {
	token, err := GetCookie(r, "session_token")
	if err != nil {
		return false, err
	}
	// compare token with database
	query := `SELECT id FROM online_status WHERE session_token = ?`
	row := db.DB.QueryRow(query, token)
	var id int
	if err := row.Scan(&id); err != nil {
		return false, err
	}
	return true, nil
}
