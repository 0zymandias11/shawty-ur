package routes

import (
	"log/slog"
	"net/http"
	"os"
	"shawty-ur/api/utils/redisUtil"
	"shawty-ur/app"
	"shawty-ur/config"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func RegisterResolveRoutes(r chi.Router, app *app.Application) {
	slog.Info("RegisterResolveRoutes called - registering /{url} route")
	r.Get("/{url}", Resolve(app))
	slog.Info("RegisterResolveRoutes completed")
}

func Resolve(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract the short URL hash from the path
		hash := chi.URLParam(req, "url")
		slog.Info("Resolving short URL", "hash", hash)

		// Look up the original URL in Redis (DB 0 - where shorten() saves URLs)
		value, err := app.RedisClient.Get(redisUtil.Ctx, hash).Result()
		if err == redis.Nil {
			slog.Warn("Short URL not found", "hash", hash)
			http.Error(w, "Short URL not found or expired", http.StatusNotFound)
			return // ✅ MUST RETURN HERE!
		} else if err != nil {
			slog.Error("Redis error while resolving URL", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return // ✅ MUST RETURN HERE!
		}

		// Connect to Redis DB 1 for analytics/counter
		rInr, err := redisUtil.New(config.RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       1,
		})
		if err != nil {
			slog.Error("Failed to connect to Redis DB 1 for analytics", "error", err)
			// Don't fail the redirect just because analytics failed
			// Just log and continue
		} else {
			defer rInr.Close()
			// Increment click counter for this URL
			_ = rInr.Incr(redisUtil.Ctx, hash)
		}

		// Redirect to the original URL
		slog.Info("Redirecting to original URL", "hash", hash, "url", value)
		http.Redirect(w, req, value, http.StatusMovedPermanently)
	}
}
