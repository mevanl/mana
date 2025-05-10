package db

import (
	"context"
	"database/sql"
	"mana/internal/models"

	"github.com/google/uuid"
)

type GuildRoleStore struct {
	DB *sql.DB
}

func NewGuildRoleStore(db *sql.DB) *GuildRoleStore {
	return &GuildRoleStore{DB: db}
}

func (guildRoleStore *GuildRoleStore) CreateGuildRole(ctx context.Context, role *models.GuildRole) error {
	insertGuildRoleSQL := `
		INSERT INTO guild_roles (id, guild_id, name, position, permissions, color, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := guildRoleStore.DB.ExecContext(ctx,
		insertGuildRoleSQL,
		role.ID,
		role.GuildID,
		role.Name,
		role.Position,
		role.Permissions,
		role.Color,
		role.CreatedAt,
	)

	return err
}

func (guildRoleStore *GuildRoleStore) GetRolesForGuild(ctx context.Context, guildID uuid.UUID) ([]*models.GuildRole, error) {
	GetGuildRolesSQL := `
		SELECT id, guild_id, name, position, permissions, color, created_at
		FROM guild_roles
		WHERE guild_id = $1
		ORDER BY position ASC
	`
	rows, err := guildRoleStore.DB.QueryContext(ctx, GetGuildRolesSQL, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.GuildRole
	for rows.Next() {
		var role models.GuildRole
		if err := rows.Scan(
			&role.ID,
			&role.GuildID,
			&role.Name,
			&role.Position,
			&role.Permissions,
			&role.Color,
			&role.CreatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, rows.Err()
}

func (guildRoleStore *GuildRoleStore) AssignRoleToMember(ctx context.Context, guildID, userID, roleID uuid.UUID) error {
	insertGuildMemberRoleSQL := `
		INSERT INTO guild_member_roles (guild_id, user_id, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := guildRoleStore.DB.ExecContext(ctx,
		insertGuildMemberRoleSQL,
		guildID,
		userID,
		roleID)

	return err
}

func (guildRoleStore *GuildRoleStore) RemoveRoleFromMember(ctx context.Context, guildID, userID, roleID uuid.UUID) error {
	deleteGuildMemberRoleSQL := `
		DELETE FROM guild_member_roles
		WHERE guild_id = $1 AND user_id = $2 AND role_id = $3
	`
	_, err := guildRoleStore.DB.ExecContext(ctx, deleteGuildMemberRoleSQL, guildID, userID, roleID)

	return err
}

func (guildRoleStore *GuildRoleStore) GetRolesForMember(ctx context.Context, guildID, userID uuid.UUID) ([]*models.GuildRole, error) {

	getGuildMemberRolesSQL := `
		SELECT gr.id, gr.guild_id, gr.name, gr.position, gr.permissions, gr.color, gr.created_at
		FROM guild_roles gr
		JOIN guild_member_roles gmr ON gr.id = gmr.role_id
		WHERE gmr.guild_id = $1 AND gmr.user_id = $2
		ORDER BY gr.position ASC
	`
	rows, err := guildRoleStore.DB.QueryContext(ctx, getGuildMemberRolesSQL, guildID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.GuildRole
	for rows.Next() {
		var role models.GuildRole
		if err := rows.Scan(
			&role.ID,
			&role.GuildID,
			&role.Name,
			&role.Position,
			&role.Permissions,
			&role.Color,
			&role.CreatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, rows.Err()
}

func (guildRoleStore *GuildRoleStore) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	deleteRoleSQL := `DELETE FROM guild_roles WHERE id = $1`
	_, err := guildRoleStore.DB.ExecContext(ctx, deleteRoleSQL, roleID)
	return err
}

func (guildRoleStore *GuildRoleStore) ReorderRoles(ctx context.Context, guildID uuid.UUID, roleIDs []uuid.UUID) error {
	tx, err := guildRoleStore.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateRoleOrderSQL := `UPDATE guild_roles SET position = $1 WHERE id = $2 AND guild_id = $3`

	for i, id := range roleIDs {
		if _, err := tx.ExecContext(ctx, updateRoleOrderSQL, i, id, guildID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (guildRoleStore *GuildRoleStore) RoleExistsByName(ctx context.Context, guildID uuid.UUID, roleName string) (bool, error) {
	selectRoleExistSQL := `
		SELECT 1 FROM guild_roles
		WHERE guild_id = $1 AND LOWER(name) = LOWER($2)
		LIMIT 1
	`
	row := guildRoleStore.DB.QueryRowContext(ctx, selectRoleExistSQL, guildID, roleName)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}
