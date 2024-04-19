package store

import (
	"database/sql"
	"fmt"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

func NewPostgresDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Check if the database connection is successful
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
