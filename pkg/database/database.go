package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return db, nil
}

func RunMigrations(databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS queries (
			id VARCHAR(255) PRIMARY KEY,
			cadastral_number VARCHAR(255) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			result BOOLEAN,
			user_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_queries_cadastral ON queries(cadastral_number)`,
		`CREATE INDEX IF NOT EXISTS idx_queries_user_id ON queries(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_queries_created_at ON queries(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_queries_status ON queries(status)`,
		`INSERT INTO users (id, username, password_hash, created_at)
		 VALUES (
			'admin_001',
			'admin',
			'$2a$10$N9qo8uLOickgx2ZMRZoMyeS7.2Y5Z1e8Z5c6W5q5k5n5v5c5n5v5c5n',
			CURRENT_TIMESTAMP
		 ) ON CONFLICT (username) DO NOTHING`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}
