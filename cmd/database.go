package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

// openDatabase opens a PostgreSQL connection using DATABASE_URL.
// Example:
//
//	DATABASE_URL=postgres://user:pass@postgres:5432/whoknows?sslmode=disable
func openDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Sensible pool defaults (tune later if needed)
	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(10)
	database.SetConnMaxLifetime(30 * time.Minute)

	// Fail fast if DB is unreachable / creds are wrong
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, err
	}

	db = database
	return database, nil
}

func closeDatabase() {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
}
