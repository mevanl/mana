package types

import (
	"context"

	"github.com/google/uuid"
)

type Client struct {
	Ctx       context.Context
	Hub       HubInterface
	Send      chan []byte
	UserID    uuid.UUID
	ChannelID uuid.UUID
}

type MessagePayload struct {
	Content string `json:"content"`
}
