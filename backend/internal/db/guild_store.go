package db

import (
	"context"
	"database/sql"
	"errors"
	"mana/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type GuildStore struct {
	DB *sql.DB
}

func NewGuildStore(db *sql.DB) *GuildStore {
	return &GuildStore{DB: db}
}

func (guildStore *GuildStore) insertGuildOnce(ctx context.Context, guild *models.Guild) error {
	insertGuildSQL := `
		INSERT INTO guilds (id, name, owner_id, invite_code, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := guildStore.DB.ExecContext(
		ctx,
		insertGuildSQL,
		guild.ID,
		guild.Name,
		guild.OwnerID,
		guild.InviteCode,
		guild.CreatedAt,
	)

	return err
}

func (guildStore *GuildStore) InsertGuild(ctx context.Context, guild *models.Guild) error {

	for i := 0; i < 3; i++ {
		err := guildStore.insertGuildOnce(ctx, guild)

		if isUniqueViolation(err, "guilds_invite_code_key") {
			guild.InviteCode = models.GenerateInviteCode()
			continue
		}

		return err
	}

	return errors.New("failed to insert guild after 3 attempts due to invite code uniqueness")
}

func (guildStore *GuildStore) DeleteGuild(ctx context.Context, guildID uuid.UUID) error {
	deleteGuildSQL := `
		DELETE FROM guilds
		WHERE id = $1
	`

	_, err := guildStore.DB.ExecContext(ctx, deleteGuildSQL, guildID)
	return err
}

func (guildStore *GuildStore) GetGuildByID(ctx context.Context, guildID uuid.UUID) (*models.Guild, error) {
	selectGuildSQL := `
		SELECT id, name, owner_id, invite_code, created_at
		FROM guilds
		WHERE id = $1
	`

	guildRow := guildStore.DB.QueryRowContext(ctx, selectGuildSQL, guildID)

	var guild models.Guild
	err := guildRow.Scan(
		&guild.ID,
		&guild.Name,
		&guild.OwnerID,
		&guild.InviteCode,
		&guild.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &guild, err
}

func (guildStore *GuildStore) GetGuildByInviteCode(ctx context.Context, inviteCode string) (*models.Guild, error) {
	selectGuildSQL := `
		SELECT id, name, owner_id, invite_code, created_at
		FROM guilds
		WHERE invite_code = $1
	`

	guildRow := guildStore.DB.QueryRowContext(ctx, selectGuildSQL, inviteCode)

	var guild models.Guild
	err := guildRow.Scan(
		&guild.ID,
		&guild.Name,
		&guild.OwnerID,
		&guild.InviteCode,
		&guild.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &guild, err
}

func (guildStore *GuildStore) AddUserToGuild(ctx context.Context, guildMember *models.GuildMember) error {
	insertUserIntoGuildSQL := `
		INSERT INTO guild_members (guild_id, user_id, joined_at)
		VALUES ($1, $2, $3)
	`

	_, err := guildStore.DB.ExecContext(
		ctx,
		insertUserIntoGuildSQL,
		guildMember.UserID,
		guildMember.GuildID,
		guildMember.JoinedAt,
	)

	return err
}

func (guildStore *GuildStore) RemoveUserFromGuild(ctx context.Context, guildID uuid.UUID, userID uuid.UUID) error {
	deleteUserFromGuildSQL := `
		DELETE FROM guild_members
		WHERE guild_id = $1 AND user_id = $2
	`
	_, err := guildStore.DB.ExecContext(ctx, deleteUserFromGuildSQL, guildID, userID)
	return err
}

func (guildStore *GuildStore) GetGuildMembers(ctx context.Context, guildID uuid.UUID) ([]*models.GuildMember, error) {
	getUsersFromGuildSQL := `
		SELECT user_id, guild_id, joined_at
		FROM guild_members
		WHERE guild_id = $1
	`

	rows, err := guildStore.DB.QueryContext(ctx, getUsersFromGuildSQL, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GuildMember

	for rows.Next() {
		var member models.GuildMember
		if err := rows.Scan(
			&member.UserID,
			&member.GuildID,
			&member.JoinedAt,
		); err != nil {
			return nil, err
		}

		members = append(members, &member)
	}

	return members, rows.Err()
}

func (guildStore *GuildStore) GetGuildsForUserID(ctx context.Context, userID uuid.UUID) ([]*models.Guild, error) {
	getUserGuildsSQL := `
		SELECT g.id, g.name, g.owner_id, g.invite_code g.created_at
		FROM guilds g
		JOIN guild_members gm ON g.id = gm.guild_id
		WHERE gm.user_id = $1
	`

	rows, err := guildStore.DB.QueryContext(ctx, getUserGuildsSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guilds []*models.Guild

	for rows.Next() {
		var guild models.Guild
		if err := rows.Scan(
			&guild.ID,
			&guild.Name,
			&guild.OwnerID,
			&guild.InviteCode,
			&guild.CreatedAt,
		); err != nil {
			return nil, err
		}
		guilds = append(guilds, &guild)
	}

	return guilds, rows.Err()
}

func (guildStore *GuildStore) CheckUserMemberOfGuild(ctx context.Context, guildID uuid.UUID, userID uuid.UUID) (bool, error) {

	return false, nil
}

func isUniqueViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.ConstraintName == constraintName
	}
	return false
}
