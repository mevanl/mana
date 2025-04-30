package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Store struct {
	db    *sql.DB
	Users *UserStore
}

func NewStore() (*Store, error) {
	argumentString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "mana"),
	)

	var err error
	var db *sql.DB

	db, err = sql.Open("postgres", argumentString)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to connect to the database: %w", err)
	}

	// Connection pool settings
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ERROR: Failed to ping to connected database: %w", err)
	}

	store := &Store{
		db:    db,
		Users: NewUserStore(db),
	}

	log.Println("Connected to PostgreSQL.")

	if err := store.createTables(); err != nil {
		return nil, fmt.Errorf("ERROR: Failed to create tables: %w", err)
	}

	return store, nil
}

func getEnv(key string, defaultFallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultFallback
}

func (store *Store) Ping() error {
	if err := store.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) createTables() error {
	// SQL Query statements
	createUserTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			activity_status TEXT DEFAULT 'offline',
			account_status TEXT DEFAULT 'active',
			created_at TIMESTAMPTZ NOT NULL
		);
	`

	_, err := store.db.Exec(createUserTableSQL)
	if err != nil {
		return err
	}
	log.Println("Users table ready.")

	log.Println("All tables ready.")
	return nil
}
