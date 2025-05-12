package router

import (
	"encoding/json"
	"log"
	"mana/internal/types"
	"mana/internal/websocket/handler"
)

type Router struct {
	Handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *Router {
	return &Router{Handler: handler}
}

func (r *Router) HandleEvent(client *types.Client, raw []byte) {
	var event types.Event

	// Parse raw incoming message into an Event struct
	if err := json.Unmarshal(raw, &event); err != nil {
		log.Printf("router: invalid event format: %v", err)
		return
	}

	switch event.Type {
	case types.EventSendMessage:

		r.Handler.HandleSendMessage(client, event.Data)

	default:
		log.Printf("router: unhandled event type: %s", event.Type)
	}
}
