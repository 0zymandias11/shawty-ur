package middleware

import (
	"net/http"
	"shawty-ur/api/metrics"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		//Track Active Connections
		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, req)

		//Get the full route pattern (e.g /api/v1/users/{id}...)
		routePattern := chi.RouteContext(req.Context()).RoutePattern()
		if routePattern == "" {
			routePattern = "unknown"
		}

		//Record metrics

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)

		metrics.HttpRequestsTotal.WithLabelValues(
			req.Method,
			routePattern,
			status,
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			req.Method,
			routePattern,
		).Observe(duration)

	})
}
