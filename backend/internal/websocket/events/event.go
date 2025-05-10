package events

import (
	"encoding/json"
	"log"
	"mana/internal/types"
)

func HandleEvent(client *types.Client, raw []byte) {
	var event types.Event

	// put raw bytes into our event struct
	if err := json.Unmarshal(raw, &event); err != nil {
		log.Printf("Invalid event format: %v", err)
		return
	}

	// handle event types
	switch event.Type {
	case types.EventSendMessage:
		handleSendMessage(client, event.Data)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}
}

func mustMarshal(data any) []byte {
	b, _ := json.Marshal(data)
	return b
}
