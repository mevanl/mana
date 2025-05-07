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
		getEnv("DB_NAME", "mana_db"),
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
	createGuildsTableSQL := `
		CREATE TABLE guilds (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`
	createGuildMembersTableSQL := `
		CREATE TABLE guild_members (
			guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (guild_id, user_id)
		);
	`

	createGuildRolesTableSQL := `
		CREATE TABLE guild_roles (
			id UUID PRIMARY KEY,
			guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			position SMALLINT NOT NULL CHECK (position BETWEEN 0 AND 255),
			permissions BIGINT NOT NULL,
			color TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`

	createGuildMemberRolesTableSQL := `
		CREATE TABLE guild_member_roles (
			guild_id UUID NOT NULL,
			user_id UUID NOT NULL,
			role_id UUID NOT NULL REFERENCES guild_roles(id) ON DELETE CASCADE,
			PRIMARY KEY (guild_id, user_id, role_id),
			FOREIGN KEY (guild_id, user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE CASCADE
		);
	`

	createGuildChannelsTableSQL := `
		CREATE TYPE channel_type AS ENUM ('text', 'voice');

		CREATE TABLE guild_channels (
			id UUID PRIMARY KEY,
			guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type channel_type NOT NULL,
			position SMALLINT NOT NULL CHECK (position BETWEEN 0 AND 255),
			topic TEXT,
			bitrate INT, 	 -- only for voice chan
			user_limit INT,  -- only for voice chan
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`

	createGuildChannelPermissionOverridesTableSQL := `
		CREATE TABLE channel_permission_overrides (
			channel_id UUID NOT NULL REFERENCES guild_channels(id) ON DELETE CASCADE,
			user_id UUID,
			role_id UUID,
			allow BIGINT NOT NULL,
			deny BIGINT NOT NULL,
			CHECK (
				(user_id IS NOT NULL AND role_id IS NULL) OR
				(user_id IS NULL AND role_id IS NOT NULL)
			),
			PRIMARY KEY (channel_id, COALESCE(user_id, role_id))
		);
	`

	var err error

	_, err = store.db.Exec(createUserTableSQL)
	if err != nil {
		return err
	}
	log.Println("Users table ready.")

	_, err = store.db.Exec(createGuildsTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guilds table ready.")

	_, err = store.db.Exec(createGuildMembersTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guild members table ready.")

	_, err = store.db.Exec(createGuildRolesTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guild roles table ready.")

	_, err = store.db.Exec(createGuildMemberRolesTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guild member roles table ready.")

	_, err = store.db.Exec(createGuildChannelsTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guild channels table ready.")

	_, err = store.db.Exec(createGuildChannelPermissionOverridesTableSQL)
	if err != nil {
		return err
	}
	log.Println("Guild channel permission overrides table ready.")

	log.Println("All tables ready.")
	return nil
}
