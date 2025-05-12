package store

import (
	"context"
	"database/sql"
	"mana/internal/models"

	"github.com/google/uuid"
)

type GuildChannelStore interface {
	Create(ctx context.Context, ch *models.GuildChannel) error
	Delete(ctx context.Context, channelID uuid.UUID) error
	Reorder(ctx context.Context, guildID uuid.UUID, channelIDs []uuid.UUID) error
	GetChannelsForGuildID(ctx context.Context, guildID uuid.UUID) ([]*models.GuildChannel, error)
	CheckNameExistsForGuildID(ctx context.Context, guildID uuid.UUID, name string) (bool, error)
}

type sqlGuildChannelStore struct {
	db *sql.DB
}

func NewGuildChannelStore(db *sql.DB) GuildChannelStore {
	return &sqlGuildChannelStore{db: db}
}

func (guildChannelStore *sqlGuildChannelStore) Create(ctx context.Context, ch *models.GuildChannel) error {

	insertChannelSQL := `
		INSERT INTO guild_channels (id, guild_id, name, type, position, topic, bitrate, user_limit, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := guildChannelStore.db.ExecContext(ctx, insertChannelSQL,
		ch.ID,
		ch.GuildID,
		ch.Name,
		ch.Type,
		ch.Position,
		ch.Topic,
		ch.Bitrate,
		ch.UserLimit,
		ch.CreatedAt,
	)
	return err
}

func (guildChannelStore *sqlGuildChannelStore) Delete(ctx context.Context, channelID uuid.UUID) error {

	deleteGuildChannelSQL := `DELETE FROM guild_channels WHERE id = $1`

	_, err := guildChannelStore.db.ExecContext(ctx, deleteGuildChannelSQL, channelID)

	return err
}

func (guildChannelStore *sqlGuildChannelStore) Reorder(ctx context.Context, guildID uuid.UUID, channelIDs []uuid.UUID) error {

	tx, err := guildChannelStore.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateGuildChannelOrderSQL := `UPDATE guild_channels SET position = $1 WHERE id = $2 AND guild_id = $3`
	for i, id := range channelIDs {
		if _, err := tx.ExecContext(ctx, updateGuildChannelOrderSQL, i, id, guildID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (guildChannelStore *sqlGuildChannelStore) GetChannelsForGuildID(ctx context.Context, guildID uuid.UUID) ([]*models.GuildChannel, error) {

	getGuildChannelsSQL := `
		SELECT id, guild_id, name, type, position, topic, bitrate, user_limit, created_at
		FROM guild_channels
		WHERE guild_id = $1
		ORDER BY position ASC
	`

	rows, err := guildChannelStore.db.QueryContext(ctx, getGuildChannelsSQL, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.GuildChannel
	for rows.Next() {
		var ch models.GuildChannel
		if err := rows.Scan(
			&ch.ID,
			&ch.GuildID,
			&ch.Name,
			&ch.Type,
			&ch.Position,
			&ch.Topic,
			&ch.Bitrate,
			&ch.UserLimit,
			&ch.CreatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}
	return channels, rows.Err()
}

func (guildChannelStore *sqlGuildChannelStore) CheckNameExistsForGuildID(ctx context.Context, guildID uuid.UUID, name string) (bool, error) {

	getChannelNameSQL := `
		SELECT 1 FROM guild_channels
		WHERE guild_id = $1 AND LOWER(name) = LOWER($2)
		LIMIT 1
	`
	var exists int
	err := guildChannelStore.db.QueryRowContext(ctx, getChannelNameSQL, guildID, name).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}
