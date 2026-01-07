package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shennawardana23/example-mcp-pub/internal/config"
	_ "github.com/lib/pq"
)

// Database represents the database connection
type Database struct {
	DB *sql.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*Database, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Health checks the database health
func (d *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("database unhealthy: %w", err)
	}

	return nil
}
