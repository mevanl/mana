package store

import (
	"context"
	"database/sql"
	"mana/internal/models"

	"github.com/google/uuid"
)

type GuildChannelOverrideStore interface {
	Create(ctx context.Context, override *models.GuildChannelPermissionOverride) error
	Delete(ctx context.Context, channelID uuid.UUID, userID, roleID *uuid.UUID) error
	GetAll(ctx context.Context, channelID uuid.UUID) ([]*models.GuildChannelPermissionOverride, error)
}

type sqlGuildChannelOverrideStore struct {
	db *sql.DB
}

func NewGuildChannelOverrideStore(db *sql.DB) GuildChannelOverrideStore {
	return &sqlGuildChannelOverrideStore{db: db}
}

func (guildChannelOverrideStore *sqlGuildChannelOverrideStore) Create(ctx context.Context, override *models.GuildChannelPermissionOverride) error {

	insertGuildChannelOverrideSQL := `
	INSERT INTO guild_channel_permission_overrides (channel_id, user_id, role_id, allow, deny)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (channel_id, user_id, role_id)
	DO UPDATE SET allow = EXCLUDED.allow, deny = EXCLUDED.deny
	`
	_, err := guildChannelOverrideStore.db.ExecContext(ctx, insertGuildChannelOverrideSQL,
		override.ChannelID,
		override.UserID,
		override.RoleID,
		override.Allow,
		override.Deny,
	)
	return err
}

func (guildChannelOverrideStore *sqlGuildChannelOverrideStore) Delete(ctx context.Context, channelID uuid.UUID, userID, roleID *uuid.UUID) error {

	deleteGuildChannelOverrideSQL := `
		DELETE FROM guild_channel_permission_overrides
		WHERE channel_id = $1 AND user_id IS NOT DISTINCT FROM $2 AND role_id IS NOT DISTINCT FROM $3
	`
	_, err := guildChannelOverrideStore.db.ExecContext(ctx, deleteGuildChannelOverrideSQL, channelID, userID, roleID)
	return err
}

func (guildChannelOverrideStore *sqlGuildChannelOverrideStore) GetAll(ctx context.Context, channelID uuid.UUID) ([]*models.GuildChannelPermissionOverride, error) {

	getGuildChannelOverrideSQL := `
		SELECT channel_id, user_id, role_id, allow, deny
		FROM guild_channel_permission_overrides
		WHERE channel_id = $1
	`
	rows, err := guildChannelOverrideStore.db.QueryContext(ctx, getGuildChannelOverrideSQL, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overrides []*models.GuildChannelPermissionOverride
	for rows.Next() {
		var o models.GuildChannelPermissionOverride
		if err := rows.Scan(&o.ChannelID, &o.UserID, &o.RoleID, &o.Allow, &o.Deny); err != nil {
			return nil, err
		}
		overrides = append(overrides, &o)
	}
	return overrides, rows.Err()
}
