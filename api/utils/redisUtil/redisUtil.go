package redisUtil

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"shawty-ur/config"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

// New creates a new Redis client
func New(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("Cannot connect to Redis", "error", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	slog.Info("Successfully connected to Redis", "addr", cfg.Addr)
	return client, nil
}
