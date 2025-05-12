package types

const (
	// General Server Permissions
	PermissionViewChannels   uint64 = 1 << 0
	PermissionManageChannels uint64 = 1 << 1
	PermissionManageRoles    uint64 = 1 << 2
	PermissionManageEmotes   uint64 = 1 << 3
	PermissionViewAudit      uint64 = 1 << 4
	PermissionViewInsights   uint64 = 1 << 5
	PermissionManageWebhooks uint64 = 1 << 6
	PermissionManageGuild    uint64 = 1 << 7

	// General Member Permissions
	PermissionChangeNickname uint64 = 1 << 10
	PermissionKickMembers    uint64 = 1 << 11
	PermissionBanMembers     uint64 = 1 << 12
	PermissionTimeoutMembers uint64 = 1 << 13

	// Text Channel Permissions
	PermissionSendMessages       uint64 = 1 << 20
	PermissionEmbedLinks         uint64 = 1 << 21
	PermissionAttachFiles        uint64 = 1 << 22
	PermissionAddReaction        uint64 = 1 << 23
	PermissionUseExternalEmotes  uint64 = 1 << 24
	PermissionMentionEveryone    uint64 = 1 << 25
	PermissionManageMessages     uint64 = 1 << 26
	PermissionReadMessageHistory uint64 = 1 << 27

	// Voice Channel Permissions
	PermissionConnect       uint64 = 1 << 30
	PermissionSpeak         uint64 = 1 << 31
	PermissionVideo         uint64 = 1 << 32
	PermissionVoiceActivity uint64 = 1 << 33
	PermissionMuteMembers   uint64 = 1 << 34
	PermissionDeafenMembers uint64 = 1 << 35
	PermissionMoveMembers   uint64 = 1 << 36

	// Admin
	PermissionAdministrator uint64 = 1 << 63
)

func HasPermission(current uint64, check uint64) bool {
	return current&check == check
}

func AddPermission(current uint64, add uint64) uint64 {
	return current | add
}

func RemovePermission(current uint64, remove uint64) uint64 {
	return current &^ remove
}
