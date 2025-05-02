package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Defines a guild
type Guild struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

// defines a member of a guild
type GuildMember struct {
	GuidID   uuid.UUID `json:"guild_id"`
	UserID   uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}


func NewGuild(name string, ownerID string) (*Guild, error) {

	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, errors.New("ownerID was invalid")
	}

	return &Guild{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   ownerUUID,
		CreatedAt: time.Now().UTC(),
	}, nil
}
