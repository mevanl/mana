package manaerror

import "errors"

var (
	ErrGuildNotFound          = errors.New("guild not found")
	ErrNotGuildOwner          = errors.New("user is not the owner of the guild")
	ErrGuildDeleteFailed      = errors.New("failed to delete guild")
	ErrBadGuildName           = errors.New("guild name must be between 2 and 100 characters")
	ErrGuildCreateFailed      = errors.New("failed to create guild")
	ErrGuildFetchFailed       = errors.New("could not get guilds")
	ErrGuildInviteNotFound    = errors.New("guild invite code not valid")
	ErrGuildJoinFailed        = errors.New("failed to join")
	ErrGuildJoinAlreadyMember = errors.New("failed to join, you are already a member")
)
