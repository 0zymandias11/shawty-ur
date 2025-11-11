package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"regexp"
	"shawty-ur/api/utils/redisUtil"
	"shawty-ur/app"
	"shawty-ur/config"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	ShortUrl       string `json:"shortUrl"`
	XRateRemaining int    `json:"rate_remaining"`
	XTimeRemaining int    `json:"time_remaining"`
}

var pattern string = "^(https?://)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}(:\\d+)?(/.*)?$"

func RegisterServiceRoutes(r chi.Router, app *app.Application) {
	// r.Get("/resolve", resolve(app))
	r.Post("/shorten", shorten(app))
}

func shorten(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("holaaa \n\n\n")
		request := new(Request)
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		}

		r2, err := redisUtil.New(config.RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       1,
		})
		if err != nil {
			slog.Info("INTERNAL SERVER ERROR ON REDIS !!")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		ip, _, err := net.SplitHostPort(req.RemoteAddr)

		if err != nil {
			http.Error(w, "Unable to parse IP", http.StatusInternalServerError)
			return
		}

		fmt.Println("Client IP:", ip)

		value, err := r2.Get(redisUtil.Ctx, ip).Result()
		if err == redis.Nil {
			_ = r2.Set(redisUtil.Ctx, ip, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
		} else {
			valInt, _ := strconv.Atoi(value)
			if valInt <= 0 {
				limit, _ := r2.TTL(redisUtil.Ctx, ip).Result()
				response := map[string]interface{}{
					"error":            "rate limit exceeded",
					"rate_limit_reset": limit / time.Nanosecond / time.Minute,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(response)
			}
		}
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(request.URL) {
			hash := uuid.New().String()
			hash = strings.ReplaceAll(hash, "-", "")[:8]

			r1, err := redisUtil.New(config.RedisConfig{
				Addr:     os.Getenv("REDIS_ADDR"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       0,
			})
			if err != nil {
				slog.Info("INTERNAL SERVER ERROR ON REDIS !!")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			val, _ := r1.Set(redisUtil.Ctx, hash, request.URL, request.Expiry*3600*time.Second).Result()
			log.Println("Printing result of set in shorten: ", val)
			resp := new(Response)
			r2.Decr(redisUtil.Ctx, ip)

			val, _ = r2.Get(redisUtil.Ctx, ip).Result()
			resp.XRateRemaining, _ = strconv.Atoi(val)
			ttl, _ := r2.TTL(redisUtil.Ctx, ip).Result()
			resp.XTimeRemaining = int(ttl / time.Nanosecond / time.Minute)
			resp.ShortUrl = os.Getenv("DOMAIN") + "/" + hash

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			response := map[string]interface{}{
				"error": "Invalid Request Url",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}

	}
}
