package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const MaxChannels uint8 = 255
const MaxRoles uint8 = 255 // 0 - 254 are free, 255 is default/'everyone' role

type ChannelType string

const (
	ChannelTypeText  ChannelType = "text"
	ChannelTypeVoice ChannelType = "voice"
)

type Guild struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OwnerID    uuid.UUID `json:"owner_id"`
	InviteCode string    `json:"invite_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type GuildMember struct {
	GuildID  uuid.UUID `json:"guild_id"`
	UserID   uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

type GuildRole struct {
	ID          uuid.UUID `json:"id"`
	GuildID     uuid.UUID `json:"guild_id"`
	Name        string    `json:"name"`
	Position    uint8     `json:"position"`
	Permissions uint64    `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	Color       string    `json:"color"`
}

type GuildMemberRole struct {
	GuildID uuid.UUID `json:"guild_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
}

type GuildChannel struct {
	ID        uuid.UUID   `json:"id"`
	GuildID   uuid.UUID   `json:"guild_id"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"type"`
	Position  uint8       `json:"position"`
	Topic     string      `json:"topic,omitempty"`
	Bitrate   *int        `json:"bitrate,omitempty"`
	UserLimit *int        `json:"user_limit,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type GuildChannelPermissionOverride struct {
	ChannelID uuid.UUID  `json:"channel_id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	RoleID    *uuid.UUID `json:"role_id,omitempty"`
	Allow     uint64     `json:"allow"`
	Deny      uint64     `json:"deny"`
}

type GuildCreateResult struct {
	Guild          *Guild
	InviteCode     string
	EveryoneRole   *GuildRole
	OwnerRole      *GuildRole
	OwnerMember    *GuildMember
	OwnerBinding   *GuildMemberRole
	GeneralChannel *GuildChannel
}

func NewGuild(name string, ownerID uuid.UUID) *Guild {
	return &Guild{
		ID:         uuid.New(),
		Name:       name,
		OwnerID:    ownerID,
		InviteCode: GenerateInviteCode(),
		CreatedAt:  time.Now().UTC(),
	}
}

func NewGuildRole(guildID uuid.UUID, name string, position uint8, perms uint64, color string) *GuildRole {
	return &GuildRole{
		ID:          uuid.New(),
		GuildID:     guildID,
		Name:        name,
		Position:    position,
		Permissions: perms,
		Color:       color,
		CreatedAt:   time.Now().UTC(),
	}
}

func newEveryoneRole(guildID uuid.UUID) *GuildRole {
	return &GuildRole{
		ID:          uuid.New(),
		GuildID:     guildID,
		Name:        "everyone",
		Position:    MaxRoles,
		Permissions: 0,
		Color:       "#000000",
		CreatedAt:   time.Now().UTC(),
	}
}

func NewGuildMember(guildID uuid.UUID, userID uuid.UUID) *GuildMember {
	return &GuildMember{
		GuildID:  guildID,
		UserID:   userID,
		JoinedAt: time.Now().UTC(),
	}
}

func NewGuildMemberRole(guildID uuid.UUID, userID uuid.UUID, roleID uuid.UUID) *GuildMemberRole {
	return &GuildMemberRole{
		GuildID: guildID,
		UserID:  userID,
		RoleID:  roleID,
	}
}

func NewGuildChannel(guildID uuid.UUID, name string, channelType ChannelType, position uint8, topic string, bitrate *int, userLimit *int) *GuildChannel {
	return &GuildChannel{
		ID:        uuid.New(),
		GuildID:   guildID,
		Name:      name,
		Type:      channelType,
		Position:  position,
		Topic:     topic,
		Bitrate:   bitrate,
		UserLimit: userLimit,
		CreatedAt: time.Now().UTC(),
	}
}

func CreateGuild(name string, ownerID uuid.UUID) *GuildCreateResult {
	guild := NewGuild(name, ownerID)
	everyoneRole := newEveryoneRole(guild.ID)
	ownerRole := NewGuildRole(guild.ID, "Owner", 0, 0xFFFFFFFFFFFFFFFF, "#000000")
	member := NewGuildMember(guild.ID, ownerID)
	generalChannel := NewGuildChannel(guild.ID, "general", ChannelTypeText, 0, "General channel", nil, nil)
	memberRole := NewGuildMemberRole(guild.ID, ownerID, ownerRole.ID)

	return &GuildCreateResult{
		Guild:          guild,
		InviteCode:     guild.InviteCode,
		EveryoneRole:   everyoneRole,
		OwnerRole:      ownerRole,
		OwnerMember:    member,
		OwnerBinding:   memberRole,
		GeneralChannel: generalChannel,
	}
}

func GenerateInviteCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	codeBytes := make([]byte, 8)

	for i := range codeBytes {
		codeBytes[i] = charset[r.Intn(len(charset))]
	}

	return string(codeBytes)
}
