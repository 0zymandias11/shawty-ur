package middleware

import (
	"context"
	"net/http"

	"shawty-ur/api/auth"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserContextKey is the key for storing user data in request context
	UserContextKey contextKey = "user"
	// IsAuthenticatedKey indicates if the user is authenticated
	IsAuthenticatedKey contextKey = "is_authenticated"
)

// RequireAuth is middleware that checks authentication but allows unauthenticated users
// Unauthenticated users will have IsAuthenticatedKey set to false in context
func RequireAuth(sessionStore *auth.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionStore.GetSession(r)

			var ctx context.Context
			if err != nil {
				// User is not authenticated - allow but mark as unauthenticated
				ctx = context.WithValue(r.Context(), IsAuthenticatedKey, false)
			} else {
				// User is authenticated - add session data to context
				ctx = context.WithValue(r.Context(), UserContextKey, session)
				ctx = context.WithValue(ctx, IsAuthenticatedKey, true)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user session data from request context
func GetUserFromContext(r *http.Request) (*auth.SessionData, bool) {
	session, ok := r.Context().Value(UserContextKey).(*auth.SessionData)
	return session, ok
}

// IsAuthenticated checks if the current request is from an authenticated user
func IsAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(IsAuthenticatedKey).(bool)
	return ok && isAuth
}
