package database

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Postgres driver
)

// NewDB initializes and returns a new database connection pool (*sql.DB).
// It takes a DSN (Data Source Name) string.
func NewDB(dsn string) (*sql.DB, error) {
	// Open a new database connection.
	// The "pgx" driver name is registered by the blank import.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters.
	// These are example values; adjust them based on your application's needs.
	db.SetMaxOpenConns(25) // Maximum number of open connections to the database.
	db.SetMaxIdleConns(25) // Maximum number of connections in the idle connection pool.
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused.

	// Verify the connection with a ping.
	// Use a context with a timeout for the ping.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		// If ping fails, close the connection pool before returning the error.
		db.Close()
		return nil, err
	}

	slog.Info("Database connection pool established successfully")
	return db, nil
}
