package authmiddleware

import "context"

// contextKey is an unexported type to prevent collisions with other packages.

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// GetUserID retrieves the authenticated user's ID from the request context.
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// GetEmail retrieves the authenticated user's email from the request context.
func GetEmail(ctx context.Context) string {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok {
		return ""
	}
	return email
}

// GetRole retrieves the authenticated user's role from the request context.
func GetRole(ctx context.Context) string {
	role, ok := ctx.Value(RoleKey).(string)
	if !ok {
		return ""
	}
	return role
}
