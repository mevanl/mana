package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ChannelType int

const TextChannel ChannelType = 0
const VoiceChannel ChannelType = 1

type Channel struct {
	ID        uuid.UUID   `json:"id"`
	GuildID   uuid.UUID   `json:"guild_id"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"type"`
	Position  int         `json:"position"` // the lower, the higher it is on the screen
	CreatedAt time.Time   `json:"created_at"`
}

func NewChannel(guildID string, name string, channelType ChannelType, position int) (*Channel, error) {
	guildUUID, err := uuid.Parse(guildID)
	if err != nil {
		return nil, errors.New("guildID was invalid")
	}

	return &Channel{
		ID:        uuid.New(),
		GuildID:   guildUUID,
		Name:      name,
		Type:      channelType,
		Position:  position,
		CreatedAt: time.Now().UTC(),
	}, nil
}
