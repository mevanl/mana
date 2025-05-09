package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	ChannelID uuid.UUID `json:"channel_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(channelID uuid.UUID, authorID uuid.UUID, content string) *Message {
	return &Message{
		ID:        uuid.New(),
		ChannelID: channelID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
}

