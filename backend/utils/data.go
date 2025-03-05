package utils

import (
	"errors"
	"real-time-forum/backend/database"
	"regexp"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func NewUUID() (uuid.UUID, error) {
	id, err := uuid.NewV4()
	return id, err
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidateUser (user database.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" || user.Age == 0 || user.Gender == "" {
		return errors.New("All fields are required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("Invalid email format")
	}
	if user.Age <= 0 {
		return errors.New("Age must be a positive number")
	}
	if user.Gender != "male" && user.Gender != "female" && user.Gender != "other" {
		return errors.New("Gender must be 'male', 'female', or 'other'")
	}
	return nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}