package db

import (
	"database/sql"
	"fmt"

	"med-predict-backend/internal/config"

	_ "github.com/lib/pq"
)

// Database wraps the SQL connection pool
type Database struct {
	DB *sql.DB
}

// Connect establishes a PostgreSQL connection
func Connect(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
