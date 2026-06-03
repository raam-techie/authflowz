package internal

import (
	"errors"
	"log"
	"new-auth-service/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// Store holds the PostgreSQL connection pool.
	Store *pgxpool.Pool

	// DB provides sqlc generated query methods.
	DB *sqlc.Store
)

// SetDB registers the shared database connection pool
// and sqlc store instance globally for application-wide access.
func SetDB(dbPool *pgxpool.Pool, sqlStore *sqlc.Store) error {
	if dbPool == nil {
		return errors.New("received nil database connection pool")
	}

	if sqlStore == nil {
		return errors.New("received nil sql store")
	}

	Store = dbPool
	DB = sqlStore

	log.Println("database dependencies initialized successfully")

	return nil
}
