package permissions

import (
	"context"
	"mana/internal/models"

	"github.com/google/uuid"
)

type PermissionStore interface {
	GetRolesForMember(ctx context.Context, guildID uuid.UUID, userID uuid.UUID) ([]*models.GuildRole, error)
	GetChannelOverrides(ctx context.Context, channelID uuid.UUID) ([]*models.GuildChannelPermissionOverride, error)
}

func ResolveBasePermissions(ctx context.Context, permissionStore PermissionStore, guildID uuid.UUID, userID uuid.UUID) (uint64, error) {
	userRoles, err := permissionStore.GetRolesForMember(ctx, guildID, userID)
	if err != nil {
		return 0, err
	}

	var userPermissions uint64
	for _, role := range userRoles {
		userPermissions |= role.Permissions
	}

	// Admin override
	if HasPermission(userPermissions, PermissionAdministrator) {
		return ^uint64(0), nil
	}

	return userPermissions, nil
}

func ResolveChannelPermissions(ctx context.Context, permissionStore PermissionStore, guildID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) (uint64, error) {
	// Get our users base permissions based on role
	basePermissions, err := ResolveBasePermissions(ctx, permissionStore, guildID, userID)
	if err != nil {
		return 0, nil
	}

	// Get our users roles
	memberRoles, err := permissionStore.GetRolesForMember(ctx, guildID, userID)
	if err != nil {
		return 0, err
	}

	// if admin, then can do whatever they want, ignore overrides
	if HasPermission(basePermissions, PermissionAdministrator) {
		return ^uint64(0), nil
	}

	// get channel overrides
	channelOverrides, err := permissionStore.GetChannelOverrides(ctx, channelID)
	if err != nil {
		return 0, err
	}

	var allow, deny uint64

	// role override
	for _, override := range channelOverrides {
		if override.RoleID != nil {
			for _, role := range memberRoles {
				if *override.RoleID == role.ID {
					deny |= override.Deny
					allow |= override.Allow
					break
				}
			}
		}
	}

	// User overrides
	for _, override := range channelOverrides {
		if override.UserID != nil && *override.UserID == userID {
			deny |= override.Deny
			allow |= override.Allow
			break // Only one user override can exist
		}
	}

	perms := basePermissions &^ deny
	perms |= allow

	return perms, nil
}
