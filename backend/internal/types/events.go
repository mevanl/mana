package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

// EventType represents all known websocket event types.
type EventType string

// Event is a shared structure used by both WebSocket and API.
type Event struct {
	Type      EventType       `json:"type"`
	ChannelID uuid.UUID       `json:"channel_id,omitempty"` // optional: not all events are channel-bound
	Data      json.RawMessage `json:"data"`
}

// Event type constants
const (
	// Message Events
	EventSendMessage    EventType = "SEND_MESSAGE"
	EventReceiveMessage EventType = "RECEIVE_MESSAGE"

	// User Events
	EventUserJoin   EventType = "USER_JOIN"
	EventUserLeave  EventType = "USER_LEAVE"
	EventUserUpdate EventType = "USER_UPDATE"

	// Guild Events
	EventGuildCreate EventType = "GUILD_CREATE"
	EventGuildUpdate EventType = "GUILD_UPDATE"
	EventGuildDelete EventType = "GUILD_DELETE"

	// Channel Events
	EventChannelCreate EventType = "CHANNEL_CREATE"
	EventChannelUpdate EventType = "CHANNEL_UPDATE"
	EventChannelDelete EventType = "CHANNEL_DELETE"

	// Typing Indicator
	EventTypingStart EventType = "TYPING_START"
	EventTypingStop  EventType = "TYPING_STOP"
)
