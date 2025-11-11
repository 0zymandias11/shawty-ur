package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"shawty-ur/api/auth"
	"shawty-ur/api/routes"
	"shawty-ur/api/utils/db"
	"shawty-ur/api/utils/redisUtil"
	"shawty-ur/app"
	"shawty-ur/config"

	"github.com/joho/godotenv"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working dir %s", err)
	}
	envPath := filepath.Join(workingDir, "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("ERR !!! Cannot load .env %s", err)
	}

	dbConfig := db.DbConfig{
		Dsn:           os.Getenv("DB_DSN"),
		Max_idle_open: 10,
		Max_open_conn: 30,
	}

	// Parse Redis DB from environment
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			redisDB = parsed
		}
	}

	redisConfig := config.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}

	cfg := config.Config{
		DbConfig:    dbConfig,
		RedisConfig: redisConfig,
		JwtSecret:   os.Getenv("JWT_SECRET"),
		Addr:        os.Getenv("ADDR"),
	}

	// Initialize PostgreSQL connection
	dbConn, err := db.New(cfg.DbConfig)
	if err != nil {
		slog.Error("Error creating db connection Client !!! ", slog.Any("err", err))
		os.Exit(1)
	}

	// Initialize Redis connection
	redisClient, err := redisUtil.New(cfg.RedisConfig)
	if err != nil {
		slog.Error("Error creating Redis connection !!! ", slog.Any("err", err))
		os.Exit(1)
	}

	// Initialize OAuth configuration
	oauthConfig := auth.NewOAuthConfig(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	// Initialize session store
	sessionStore := auth.NewSessionStore(os.Getenv("SESSION_KEY"))

	application := &app.Application{
		Config:       cfg,
		DbConnector:  dbConn,
		RedisClient:  redisClient,
		OAuthConfig:  oauthConfig,
		SessionStore: sessionStore,
	}

	// Register all route handlers
	// Adding new routes is as simple as adding new RegisterXRoutes functions here
	application.RegisterRoutes(
		routes.RegisterHealthRoutes,
		routes.RegisterUserRoutes,
		routes.RegisterServiceRoutes,
		routes.RegisterAuthRoutes,
	)
	application.RegisterSoloRoutes(
		routes.RegisterResolveRoutes,
	)

	mux := application.Mount()
	if err := application.Run(mux); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
