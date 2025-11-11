package config

import "shawty-ur/api/utils/db"

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Config holds the application configuration
type Config struct {
	Addr        string
	JwtSecret   string
	DbConfig    db.DbConfig
	RedisConfig RedisConfig
}
