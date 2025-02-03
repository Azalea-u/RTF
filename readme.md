***File hiarchy in the repo
```bash
/forum-app
│
├── /backend
│   ├── /api
│   │   ├── handlers.go          # HTTP handlers (REST endpoints)
│   │   ├── routes.go            # Defines API routes and WebSocket setup
│   │   ├── middleware.go        # CORS, auth middleware, rate limiting
│   │   └── websocket.go         # WebSocket connection logic (using Gorilla)
│   ├── /config
│   │   └── config.go            # Load env/config (e.g., port, JWT secrets)
│   ├── /db
│   │   ├── database.go          # SQLite3 connection and setup
│   │   ├── models.go            # Structs for User, Post, Message, etc.
│   │   └── migrations           # SQL schema migration files
│   ├── /services
│   │   ├── auth_service.go      # User auth (bcrypt for passwords, JWT tokens)
│   │   ├── forum_service.go     # CRUD for posts, threads, comments
│   │   └── message_service.go   # Real-time messaging (WebSocket + SQLite storage)
│   ├── /utils
│   │   ├── logger.go            # Custom logger
│   │   ├── helpers.go           # UUID generation, sanitize inputs, etc.
│   │   └── jwt.go               # JWT token creation/validation
│   ├── main.go                  # Server entry point (starts HTTP/WS server)
│   └── go.mod                   # Go dependencies (incl. Gorilla, sqlite3, bcrypt, uuid)
│
├── /frontend
│   ├── /css
│   │   └── styles.css           # Styling for forum and chat
│   ├── /js
│   │   ├── app.js               # SPA routing and core logic
│   │   ├── auth.js              # Login/register UI and API calls
│   │   ├── forum.js             # Render posts, handle voting/comments
│   │   ├── messaging.js         # Chat UI and WebSocket client logic
│   │   └── websocket.js         # WebSocket connection manager (Gorilla client)
│   ├── /assets
│   │   └── ...                  # Icons, images, etc.
│   ├── index.html               # Single HTML entry point
│   └── 404.html                 # Fallback page
│
├── .gitignore                   # Ignore binaries, .env, node_modules, etc.
├── README.md                    # Setup instructions
└── Dockerfile                   # Optional for containerization

```