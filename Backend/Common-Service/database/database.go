// Package database opens and verifies a PostgreSQL connection pool (pgx).
// It holds NO business logic — just connectivity that services share.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kidkon/ecommerce/common/config"
)

// Config holds PostgreSQL connection settings.
type Config struct {
	Host, Port, User, Password, Name, SSLMode string
}

// dsn builds a key=value connection string (safer than a URL: no need to
// URL-encode special characters in the password).
func (c Config) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Connect opens a pgx pool and Pings once to confirm the DB is reachable.
// Returns an error (instead of panicking) so the caller decides what to do.
func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	// verify the DB actually answers, with a short timeout so startup can't hang
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}

// ConnectFromEnv reads the standard DB_* env vars and connects.
func ConnectFromEnv(ctx context.Context) (*pgxpool.Pool, error) {
	return Connect(ctx, Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.Get("DB_PORT", "5432"),
		User:     config.Get("DB_USER", "postgres"),
		Password: config.Get("DB_PASSWORD", ""),
		Name:     config.Get("DB_NAME", "postgres"),
		SSLMode:  config.Get("DB_SSLMODE", "disable"),
	})
}
