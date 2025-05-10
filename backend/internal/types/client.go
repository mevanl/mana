// defines share client types for websocket and event
package types

import "github.com/google/uuid"

type Client struct {
	Hub       HubInterface
	Send      chan []byte
	UserID    uuid.UUID
	ChannelID uuid.UUID
}

