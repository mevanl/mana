package db

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTables(db *sql.DB) error {

	// create channel_type enum
	_, err := db.Exec(`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'channel_type') THEN
				CREATE TYPE channel_type AS ENUM ('text', 'voice');
			END IF;
		END$$;`)
	if err != nil {
		return fmt.Errorf("creating enum channel_type: %w", err)
	}

	statements := []struct {
		name string
		sql  string
	}{
		{"users", `
			CREATE TABLE IF NOT EXISTS users (
				id UUID PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				email TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				activity_status TEXT DEFAULT 'offline',
				account_status TEXT DEFAULT 'active',
				created_at TIMESTAMPTZ NOT NULL
			);`},

		{"guilds", `
			CREATE TABLE IF NOT EXISTS guilds (
				id UUID PRIMARY KEY,
				name TEXT NOT NULL,
				owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				invite_code TEXT NOT NULL UNIQUE,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);`},

		{"guild_members", `
			CREATE TABLE IF NOT EXISTS guild_members (
				guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
				user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				PRIMARY KEY (guild_id, user_id)
			);`},

		{"guild_roles", `
			CREATE TABLE IF NOT EXISTS guild_roles (
				id UUID PRIMARY KEY,
				guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
				name TEXT NOT NULL,
				position SMALLINT NOT NULL CHECK (position BETWEEN 0 AND 255),
				permissions BIGINT NOT NULL,
				color TEXT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);`},

		{"guild_member_roles", `
			CREATE TABLE IF NOT EXISTS guild_member_roles (
				guild_id UUID NOT NULL,
				user_id UUID NOT NULL,
				role_id UUID NOT NULL REFERENCES guild_roles(id) ON DELETE CASCADE,
				PRIMARY KEY (guild_id, user_id, role_id),
				FOREIGN KEY (guild_id, user_id) REFERENCES guild_members(guild_id, user_id) ON DELETE CASCADE
			);`},

		{"guild_channels", `
			CREATE TABLE IF NOT EXISTS guild_channels (
				id UUID PRIMARY KEY,
				guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
				name TEXT NOT NULL,
				type channel_type NOT NULL,
				position SMALLINT NOT NULL CHECK (position BETWEEN 0 AND 255),
				topic TEXT,
				bitrate INT, 	 -- only for voice chan
				user_limit INT,  -- only for voice chan
				created_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);`},

		{"guild_channel_permission_overrides", `
			CREATE TABLE IF NOT EXISTS guild_channel_permission_overrides (
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
			);`},

		{"messages", `
			CREATE TABLE IF NOT EXISTS messages (
				id UUID PRIMARY KEY,
				channel_id UUID REFERENCES guild_channels(id) ON DELETE CASCADE,
				author_id UUID REFERENCES users(id) ON DELETE CASCADE,
				content TEXT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);
		`},
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt.sql); err != nil {
			return fmt.Errorf("error creating table %s: %w", stmt.name, err)
		}
		log.Printf("Table %s ready.", stmt.name)
	}

	return nil
}
