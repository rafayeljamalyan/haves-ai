// Package db provides the Postgres connection pool used across the service.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool. It is safe for concurrent use.
type DB struct {
	Pool *pgxpool.Pool
}

// New opens a connection pool for the given DSN and verifies connectivity
// before returning. The caller owns the returned *DB and must Close it.
func New(ctx context.Context, dsn string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db: creating pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: pinging database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Ping verifies the database is reachable. Used by readiness checks.
func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Close releases all connections in the pool.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
