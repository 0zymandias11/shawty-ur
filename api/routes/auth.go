package routes

import (
	"log/slog"
	"net/http"

	"shawty-ur/api/auth"
	"shawty-ur/api/helper"
	"shawty-ur/api/utils"
	"shawty-ur/app"

	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes registers OAuth authentication routes
func RegisterAuthRoutes(r chi.Router, application *app.Application) {
	r.Get("/auth/google", googleLoginHandler(application))
	r.Get("/auth/google/callback", googleCallbackHandler(application))
	r.Post("/auth/logout", logoutHandler(application))
	r.Get("/auth/me", meHandler(application))
}

// googleLoginHandler initiates Google OAuth flow
func googleLoginHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate state token for CSRF protection
		state, err := auth.GenerateStateToken()
		if err != nil {
			slog.Error("Failed to generate state token", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to initiate OAuth"})
			return
		}

		// Save state to session
		if err := application.SessionStore.SaveState(w, r, state); err != nil {
			slog.Error("Failed to save state to session", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to initiate OAuth"})
			return
		}

		// Redirect to Google OAuth consent page
		url := application.OAuthConfig.GoogleConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// googleCallbackHandler handles the OAuth callback from Google
func googleCallbackHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate state token
		state := r.URL.Query().Get("state")
		valid, err := application.SessionStore.ValidateState(r, state)
		if err != nil || !valid {
			slog.Error("Invalid state token", "error", err)
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid state token"})
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			slog.Error("No authorization code provided")
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "No authorization code"})
			return
		}

		// Exchange code for user info
		googleUser, err := application.OAuthConfig.GetGoogleUserInfo(r.Context(), code)
		if err != nil {
			slog.Error("Failed to get Google user info", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to authenticate with Google"})
			return
		}

		// Find or create user
		userStore := helper.NewUserStore(application.DbConnector)
		user, err := userStore.FindOrCreateOAuthUser(r.Context(), googleUser)
		if err != nil {
			slog.Error("Failed to find or create OAuth user", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
			return
		}

		// Create session
		sessionData := auth.SessionData{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			Provider: "google",
		}

		if err := application.SessionStore.SaveSession(w, r, sessionData); err != nil {
			slog.Error("Failed to save session", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
			return
		}

		// Return success response with user info
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Successfully authenticated",
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"provider": "google",
			},
		})
	}
}

// logoutHandler destroys the user session
func logoutHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := application.SessionStore.DestroySession(w, r); err != nil {
			slog.Error("Failed to destroy session", "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
	}
}

// meHandler returns the current authenticated user's information
func meHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := application.SessionStore.GetSession(r)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"user": map[string]interface{}{
				"id":       session.UserID,
				"username": session.Username,
				"email":    session.Email,
				"provider": session.Provider,
			},
		})
	}
}
