package routes

import (
	"net/http"

	"shawty-ur/app"

	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(r chi.Router, application *app.Application) {
	r.Get("/health", healthHandler(application))
}

func healthHandler(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
