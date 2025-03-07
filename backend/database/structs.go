package database

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Age       int       `db:"age" json:"age"`
	Token     []byte    `db:"token" json:"token"`
	Gender    string    `db:"gender" json:"gender"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Online    bool      `json:"online"`
}

// Post represents a post created by a user.
type Post struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Category  string    `db:"category" json:"category"`
	UserID    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Comment represents a comment made on a post.
type Comment struct {
	ID        int       `db:"id" json:"id"`
	Content   string    `db:"content" json:"content"`
	UserID    int       `db:"user_id" json:"user_id"`
	PostID    int       `db:"post_id" json:"post_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Message represents a message sent between users.
type Message struct {
	ID         int       `db:"id" json:"id"`
	SenderID   int       `db:"sender_id" json:"sender_id"`
	ReceiverID int       `db:"receiver_id" json:"receiver_id"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
