package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type DbConfig struct {
	Dsn           string
	Max_idle_open int
	Max_open_conn int
}

func New(dbConfig DbConfig) (*sql.DB, error) {
	fmt.Printf("in side NEw \n")
	if dbConfig.Dsn == "" {
		log.Fatalf("Cannot Connect to DB, DSN not present")
	}

	dsn := dbConfig.Dsn
	if !strings.Contains(dsn, "sslmode=") {
		if dsn[len(dsn)-1] != '?' {
			dsn += "?"
		}
		dsn += "sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Error opening db connection !\n\n")
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.Max_open_conn)
	db.SetMaxIdleConns(dbConfig.Max_idle_open)
	duration, err := time.ParseDuration("10m")
	if err != nil {
		slog.Error("Error setting Db timing limit\n")
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Cannot ping DB !! %s %s", dsn, err)
		return nil, err
	}
	return db, nil
}
