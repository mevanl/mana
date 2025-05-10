// defines share event types for websocket and event
package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event struct {
	Type      string          `json:"type"`
	ChannelID uuid.UUID       `json:"channel_id"`
	Data      json.RawMessage `json:"data"`
}

const (
	EventSendMessage    = "SEND_MESSAGE"
	EventReceiveMessage = "RECEIVE_MESSAGE"
)
