package utils

import "context"

type contextKey string

const userIDKey contextKey = "userID"

// SetUserIDInContext adds the user ID to the context
func SetUserIDInContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}