package handler

import (
	"encoding/json"
	"log"
	"mana/internal/services"
	"mana/internal/types"
	"mana/internal/websocket/util"
)

// HandleSendMessage handles the SEND_MESSAGE event from a client.
func (h *Handler) HandleSendMessage(client *types.Client, raw json.RawMessage) {
	var payload types.MessagePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("handler: invalid SEND_MESSAGE payload: %v", err)
		return
	}

	msg, err := services.SendMessage(client.Ctx, h.store, client.UserID, client.ChannelID, payload.Content)
	if err != nil {
		log.Printf("handler: failed to send message: %v", err)
		return
	}

	client.Hub.BroadcastMessage(types.Event{
		Type:      types.EventReceiveMessage,
		ChannelID: client.ChannelID,
		Data:      util.MustMarshal(msg),
	})
}
