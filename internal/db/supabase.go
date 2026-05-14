package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB is the package-level database connection pool.
var DB *sql.DB

// Connect opens a connection to the Supabase/Postgres database using the
// provided connection string and verifies the connection with a ping.
func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required but was not set")
	}

	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(5)

	DB = database
	log.Println("Database connection established successfully")
	return database, nil
}

// Close closes the database connection pool.
func Close(database *sql.DB) {
	if database != nil {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
