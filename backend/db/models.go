package db

import (
	"time"
)

// User represents the users table in the database.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}

// Post represents the posts table in the database.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents the comments table in the database.
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Like represents the likes table in the database.
type Like struct {
	ID        int       `json:"id"`
	PostID    *int      `json:"post_id"`
	CommentID *int      `json:"comment_id"`
	UserID    int       `json:"user_id"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Conversation represents the conversations table in the database.
type Conversation struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents the messages table in the database.
type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Content        string    `json:"content"`
	Read           bool      `json:"read"`
	CreatedAt      time.Time `json:"created_at"`
}

// OnlineStatus represents the online_status table in the database.
type OnlineStatus struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Online    bool      `json:"online"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}
