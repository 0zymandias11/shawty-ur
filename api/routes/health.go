package routes

import (
	"net/http"

	"shawty-ur/app"

	"github.com/go-chi/chi/v5"
)

// RegisterHealthRoutes registers all health-related routes
// This function is called by the Application during route setup
func RegisterHealthRoutes(r chi.Router, application *app.Application) {
	r.Get("/health", healthHandler(application))
}

// healthHandler returns the health check handler
// It's unexported (private) to keep the API clean - only RegisterHealthRoutes is exposed
func healthHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add database ping check if needed
		// if err := application.DbConnector.Ping(); err != nil {
		//     w.WriteHeader(http.StatusServiceUnavailable)
		//     return
		// }

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
