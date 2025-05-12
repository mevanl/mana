package services

import (
	"context"
	manaerror "mana/internal/errors"
	"mana/internal/models"
	"mana/internal/store"
	"strings"

	"github.com/google/uuid"
)

func CreateGuild(ctx context.Context, store *store.Store, name string, ownerID uuid.UUID) (*models.Guild, error) {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return nil, manaerror.ErrBadGuildName
	}

	guild := models.NewGuild(name, ownerID)

	if err := store.Guilds.Create(ctx, guild); err != nil {
		return nil, manaerror.ErrGuildCreateFailed
	}

	return guild, nil
}

func DeleteGuild(ctx context.Context, store *store.Store, guildID uuid.UUID, userID uuid.UUID) error {

	// get guild
	guild, err := store.Guilds.FindByID(ctx, guildID)
	if err != nil {
		return manaerror.ErrGuildNotFound
	}

	// check they are owner
	if guild.OwnerID != userID {
		return manaerror.ErrNotGuildOwner
	}

	// delete
	if err := store.Guilds.Delete(ctx, guildID); err != nil {
		return manaerror.ErrGuildDeleteFailed
	}

	return nil

}

func GetGuild(ctx context.Context, store *store.Store, guildID uuid.UUID) (*models.Guild, error) {
	guild, err := store.Guilds.FindByID(ctx, guildID)
	if err != nil {
		return nil, manaerror.ErrGuildNotFound
	}

	return guild, nil
}

func GetUserGuilds(ctx context.Context, store *store.Store, userID uuid.UUID) ([]*models.Guild, error) {
	guilds, err := store.Guilds.FindGuildsForUserID(ctx, userID)
	if err != nil {
		return nil, manaerror.ErrGuildFetchFailed
	}

	return guilds, nil
}

func JoinGuildByInvite(ctx context.Context, store *store.Store, userID uuid.UUID, inviteCode string) (*models.Guild, error) {
	guild, err := store.Guilds.FindByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, manaerror.ErrGuildInviteNotFound
	}

	exists, err := store.Guilds.CheckMembershipByUserID(ctx, guild.ID, userID)
	if err != nil {
		return nil, manaerror.ErrGuildJoinFailed
	}
	if exists {
		return nil, manaerror.ErrGuildJoinAlreadyMember
	}

	member := models.NewGuildMember(guild.ID, userID)
	if err := store.Guilds.AddMember(ctx, member); err != nil {
		return nil, manaerror.ErrGuildJoinFailed
	}

	// TODO: Add user to @everyone role if required

	return guild, nil
}
