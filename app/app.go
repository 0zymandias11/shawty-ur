package app

import (
	"database/sql"
	"log/slog"
	"net/http"

	"shawty-ur/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

// RouteRegistrar is a function that registers routes on a chi.Router
// This pattern allows routes to be defined in separate packages and injected
type RouteRegistrar func(r chi.Router, app *Application)

// Application holds the application state and dependencies
type Application struct {
	Config              config.Config
	DbConnector         *sql.DB
	RedisClient         *redis.Client
	routeRegistrars     []RouteRegistrar
	soloRouteRegistrars []RouteRegistrar
}

func (app *Application) RegisterRoutes(registrars ...RouteRegistrar) {
	app.routeRegistrars = append(app.routeRegistrars, registrars...)
}
func (app *Application) RegisterSoloRoutes(registrars ...RouteRegistrar) {
	app.soloRouteRegistrars = append(app.soloRouteRegistrars, registrars...)
}

// Mount sets up all the routes and middleware for the application
func (app *Application) Mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check at root level (must come before catch-all routes)
	r.Get("/health", app.healthCheckHandler)

	// Register versioned API routes
	r.Route("/api/v1", func(r chi.Router) {
		for _, registrar := range app.routeRegistrars {
			registrar(r, app)
		}
	})

	// Register solo routes (like /:url) directly at root level
	// These must come AFTER specific routes like /health to avoid conflicts
	slog.Info("Registering solo routes", "count", len(app.soloRouteRegistrars))
	for _, registrar := range app.soloRouteRegistrars {
		registrar(r, app)
	}

	// Print all registered routes for debugging
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info("Registered route", "method", method, "route", route)
		return nil
	})

	return r
}

// healthCheckHandler provides a quick health check at root level
func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Run starts the HTTP server
func (app *Application) Run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:    app.Config.Addr,
		Handler: mux,
	}

	slog.Info("server started at ", app.Config.Addr, "default")
	return srv.ListenAndServe()
}
