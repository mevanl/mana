package events

import (
	"encoding/json"
	"log"
	"mana/internal/types"
)

func handleSendMessage(client *types.Client, raw json.RawMessage) {
	var payload types.MessagePayload

	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("Invalid SEND_MESSAGE payload: %v", err)
		return
	}

	// Db logic here ?

	// Echo to others as RECEIVE_MESSAGE event
	// client.Hub.Broadcast <- types.Event{
	// 	Type:      types.EventReceiveMessage,
	// 	ChannelID: client.ChannelID,
	// 	Data: mustMarshal(types.MessagePayload{
	// 		Content: payload.Content,
	// 	}),
	// }

	client.Hub.BroadcastMessage(types.Event{
		Type:      types.EventReceiveMessage,
		ChannelID: client.ChannelID,
		Data:      mustMarshal(payload),
	})
}
