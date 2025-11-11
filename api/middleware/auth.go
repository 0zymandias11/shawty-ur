package middleware

import (
	"context"
	"net/http"

	"shawty-ur/api/auth"
	"shawty-ur/api/utils"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserContextKey is the key for storing user data in request context
	UserContextKey contextKey = "user"
)

// RequireAuth is middleware that requires authentication
func RequireAuth(sessionStore *auth.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionStore.GetSession(r)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
				return
			}

			// Add session data to request context
			ctx := context.WithValue(r.Context(), UserContextKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user session data from request context
func GetUserFromContext(r *http.Request) (*auth.SessionData, bool) {
	session, ok := r.Context().Value(UserContextKey).(*auth.SessionData)
	return session, ok
}
