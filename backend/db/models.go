package db

import "time"

// User represents a user in the database.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}

// Post represents a post created by a user.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Like represents a like on a post or a comment.
type Like struct {
	ID        int       `json:"id"`
	PostID    *int      `json:"post_id,omitempty"`
	CommentID *int      `json:"comment_id,omitempty"`
	UserID    int       `json:"user_id"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Conversation represents a private conversation between two users.
type Conversation struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents a message exchanged in a conversation.
type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Content        string    `json:"content"`
	Read           bool      `json:"read"`
	CreatedAt      time.Time `json:"created_at"`
}

// OnlineStatus represents the online status of a user.
type OnlineStatus struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Online    bool      `json:"online"`
	LastSeen  time.Time `json:"last_seen"`
}
