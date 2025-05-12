package services

import (
	"context"
	manaerror "mana/internal/errors"
	"mana/internal/models"
	"mana/internal/store"
	"time"

	"github.com/google/uuid"
)

func SendMessage(ctx context.Context, store *store.Store, userID uuid.UUID, channelID uuid.UUID, content string) (*models.Message, error) {
	if content == "" {
		return nil, manaerror.ErrMessageEmpty
	}

	// TODO: Validate user can send message in the channel (permissions)

	msg := models.NewMessage(channelID, userID, content)
	if err := store.Messages.Create(ctx, msg); err != nil {
		return nil, manaerror.ErrMessageSendFailed
	}

	return msg, nil
}

func GetChannelMessages(ctx context.Context, store *store.Store, channelID uuid.UUID, limit int, before *time.Time) ([]*models.Message, error) {

	messages, err := store.Messages.GetChannelMessages(ctx, channelID, limit, before)
	if err != nil {
		return nil, manaerror.ErrMessageFetchFailed
	}

	return messages, nil
}
