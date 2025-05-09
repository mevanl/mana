package db

import (
	"context"
	"database/sql"
	"mana/internal/models"

	"github.com/google/uuid"
)

type MessageStore struct {
	DB *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{DB: db}
}

func (messageStore *MessageStore) InsertMessage(ctx context.Context, message *models.Message) error {
	insertMessageSQL := `
		INSERT INTO messages (id, channel_id, author_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := messageStore.DB.ExecContext(ctx, insertMessageSQL,
		message.ID, message.ChannelID, message.AuthorID, message.Content, message.CreatedAt,
	)
	return err
}

func (messageStore *MessageStore) GetMessagesByChannel(ctx context.Context, channelID uuid.UUID, limit int) ([]*models.Message, error) {
	selectMessagesFromChannelSQL := `
		SELECT id, channel_id, author_id, content, created_at
		FROM messages
		WHERE channel_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := messageStore.DB.QueryContext(ctx, selectMessagesFromChannelSQL, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.ChannelID,
			&msg.AuthorID,
			&msg.Content,
			&msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
