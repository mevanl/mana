package store

import "database/sql"

type Store struct {
	db                    *sql.DB
	Users                 UserStore
	Guilds                GuildStore
	GuildRoles            GuildRoleStore
	GuildChannels         GuildChannelStore
	GuildChannelOverrides GuildChannelOverrideStore
	Messages              MessageStore
}

func New(db *sql.DB) *Store {
	return &Store{
		db:                    db,
		Users:                 NewUserStore(db),
		Guilds:                NewGuildStore(db),
		GuildRoles:            NewGuildRoleStore(db),
		GuildChannels:         NewGuildChannelStore(db),
		GuildChannelOverrides: NewGuildChannelOverrideStore(db),
		Messages:              NewMessageStore(db),
	}
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) Ping() error {
	if err := store.db.Ping(); err != nil {
		return err
	}
	return nil
}
