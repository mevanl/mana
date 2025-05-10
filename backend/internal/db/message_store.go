package db

import (
	"context"
	"database/sql"
	"fmt"
	"mana/internal/models"
	"time"

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

func (messageStore *MessageStore) GetMessagesByChannel(ctx context.Context, channelID uuid.UUID, limit int, before *time.Time) ([]*models.Message, error) {
	selectMessagesFromChannelSQL := `
		SELECT id, channel_id, author_id, content, created_at
		FROM messages
		WHERE channel_id = $1
	`

	// create slice of all our arguments
	args := []interface{}{channelID}
	paramIndex := 2 // keep track of # of positional args for our query

	// check if before (for pagination)
	if before != nil {
		selectMessagesFromChannelSQL += fmt.Sprintf(" AND created_at < $%d", paramIndex)
		args = append(args, *before)
		paramIndex++
	}

	// add our limit
	selectMessagesFromChannelSQL += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", paramIndex)
	args = append(args, limit)

	// execute
	rows, err := messageStore.DB.QueryContext(ctx, selectMessagesFromChannelSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// convert rows -> messages
	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.ChannelID,
			&msg.AuthorID,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	// reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
