package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// defines a role in a guild
type GuildRole struct {
	ID          uuid.UUID `json:"id"`
	GuildID     uuid.UUID `json:"guild_id"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`    // Lower number = highest position
	Permissions int64     `json:"permissions"` // Bitfield representing different permissions
	CreatedAt   time.Time `json:"created_at"`
}

// defines members who are in a role for a given guild
type GuildMemberRole struct {
	GuildID uuid.UUID `json:"guild_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
}

func NewRole(guildID string, name string, position int) (*GuildRole, error) {

	guildUUID, err := uuid.Parse(guildID)
	if err != nil {
		return nil, errors.New("guildID was invalid")
	}

	return &GuildRole{
		ID:          uuid.New(),
		GuildID:     guildUUID,
		Name:        name,
		Position:    position,
		Permissions: 0,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
