package auth

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName     = "shawty_session"
	SessionUserID   = "user_id"
	SessionUsername = "username"
	SessionEmail    = "email"
	SessionProvider = "provider"
	SessionState    = "oauth_state"
)

// SessionStore manages user sessions
type SessionStore struct {
	store *sessions.CookieStore
}

// SessionData holds user session information
type SessionData struct {
	UserID   int64
	Username string
	Email    string
	Provider string
}

func init() {
	// Register SessionData for gob encoding
	gob.Register(&SessionData{})
}

// NewSessionStore creates a new session store
func NewSessionStore(sessionKey string) *SessionStore {
	store := sessions.NewCookieStore([]byte(sessionKey))

	// Configure session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   true,  // HTTPS only
		SameSite: http.SameSiteLaxMode,
	}

	return &SessionStore{store: store}
}

// SaveSession saves user data to session
func (s *SessionStore) SaveSession(w http.ResponseWriter, r *http.Request, data SessionData) error {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Values[SessionUserID] = data.UserID
	session.Values[SessionUsername] = data.Username
	session.Values[SessionEmail] = data.Email
	session.Values[SessionProvider] = data.Provider

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetSession retrieves user data from session
func (s *SessionStore) GetSession(r *http.Request) (*SessionData, error) {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	userID, ok := session.Values[SessionUserID].(int64)
	if !ok {
		return nil, fmt.Errorf("session not found or invalid")
	}

	username, _ := session.Values[SessionUsername].(string)
	email, _ := session.Values[SessionEmail].(string)
	provider, _ := session.Values[SessionProvider].(string)

	return &SessionData{
		UserID:   userID,
		Username: username,
		Email:    email,
		Provider: provider,
	}, nil
}

// DestroySession removes the session
func (s *SessionStore) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}

	return nil
}

// SaveState saves OAuth state to session for CSRF protection
func (s *SessionStore) SaveState(w http.ResponseWriter, r *http.Request, state string) error {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Values[SessionState] = state

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// ValidateState validates OAuth state for CSRF protection
func (s *SessionStore) ValidateState(r *http.Request, state string) (bool, error) {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}

	savedState, ok := session.Values[SessionState].(string)
	if !ok || savedState == "" {
		return false, fmt.Errorf("no state found in session")
	}

	return savedState == state, nil
}

// IsAuthenticated checks if the user is authenticated
func (s *SessionStore) IsAuthenticated(r *http.Request) bool {
	_, err := s.GetSession(r)
	return err == nil
}
