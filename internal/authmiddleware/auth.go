package authmiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/devharnold/smart-reconcile/internal/auth"
)

type AuthUser struct {
	UserID string
	Email  string
	Role   string
}

type contextKey string

const UserContextKey contextKey = "auth_user"

type ErrorResponse struct {
	Error string `json:"error"`
}

func JWTAuth(jwtSvc auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing Authorization header")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "invalid Authorization header")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if token == "" {
				writeError(w, http.StatusUnauthorized, "empty token")
				return
			}

			claims, err := jwtSvc.ValidateToken(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			user := &AuthUser{
				UserID: claims.UserID,
				Email:  claims.Email,
				Role:   claims.Role,
			}

			ctx := context.WithValue(
				r.Context(),
				UserContextKey,
				user,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func GetUser(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*AuthUser)
	return user, ok
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
