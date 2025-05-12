package services

import (
	"context"
	manaerror "mana/internal/errors"
	"mana/internal/models"
	"mana/internal/store"

	"github.com/google/uuid"
)

func GetChannels(ctx context.Context, store *store.Store, guildID uuid.UUID) ([]*models.GuildChannel, error) {
	channels, err := store.GuildChannels.GetChannelsForGuildID(ctx, guildID)
	if err != nil {
		return nil, manaerror.ErrChannelFetchFailed
	}

	return channels, nil
}
