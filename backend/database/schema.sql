-- User Table --
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    age INTEGER CHECK(age > 0) NOT NULL,
    token BLOB NOT NULL,
    gender TEXT CHECK(gender IN ('male', 'female', 'other')) NOT NULL
)

-- Post Table --
CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id)
)

-- Comment Table --
CREATE TABLE IF NOT EXISTS comment (
    id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(post_id) REFERENCES post(id)
)

-- message table --
CREATE TABLE IF NOT EXISTS message (
    id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY(sender_id) REFERENCES user(id),
    FOREIGN KEY(receiver_id) REFERENCES user(id)
)