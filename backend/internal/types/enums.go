package types

// ChannelType is the type of a guild channel.
type ChannelType string

const (
	ChannelTypeText  ChannelType = "text"
	ChannelTypeVoice ChannelType = "voice"
)

// PresenceStatus represents user presence
type PresenceStatus string

const (
	PresenceOnline       PresenceStatus = "online"
	PresenceIdle         PresenceStatus = "idle"
	PresenceDoNotDisturb PresenceStatus = "dnd"
	PresenceOffline      PresenceStatus = "offline"
)

// MessageType differentiates messages (mainly for formatting)
type MessageType string

const (
	MessageTypeDefault MessageType = "default"
	MessageTypeSystem  MessageType = "system"
	MessageTypeJoin    MessageType = "user_join"
	MessageTypeLeave   MessageType = "user_leave"
)
