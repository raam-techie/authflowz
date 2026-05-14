package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transaction
type Store struct {
	*Queries
	CoonPool *pgxpool.Pool // Write DB connection pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		CoonPool: connPool,
		Queries:  New(connPool),
	}
}

// NewTx starts a new database transaction using the write connection pool
func (store *Store) NewTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := store.CoonPool.Begin(ctx) // Start a transaction using the write connection pool
	if err != nil {
		return nil, err
	}
	return tx, nil
}
