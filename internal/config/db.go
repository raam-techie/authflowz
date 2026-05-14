package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewStore creates and initializes a PostgreSQL connection pool.
//
// It configures connection pool settings such as:
//   - Maximum open connections
//   - Minimum idle connections
//   - Connection lifetime
//   - Idle timeout
//   - Health check interval
//
// The function also validates the database connection
// by performing a ping before returning the pool
func NewStore(dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	// Pool configurations
	config.MaxConns = 10
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	// Verify database connection
	if err := connPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Println("database connection established successfully")

	return connPool, nil
}
