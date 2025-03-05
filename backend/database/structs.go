package database

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Age       int       `db:"age"`
	Token     []byte    `db:"token"`
	Gender    string    `db:"gender"`
	CreatedAt time.Time `db:"created_at"`
}

type Post struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Category  string    `db:"category"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Comment struct {
	ID        int       `db:"id"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id"`
	PostID    int       `db:"post_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Message struct {
	ID         int       `db:"id"`
	SenderID   int       `db:"sender_id"`
	ReceiverID int       `db:"receiver_id"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
}
