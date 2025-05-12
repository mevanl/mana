package websocket

import (
	"sync"

	"mana/internal/types"

	"github.com/google/uuid"
)

type Hub struct {
	mu sync.RWMutex

	// Maps channel IDs to clients currently connected
	channels map[uuid.UUID]map[*types.Client]bool

	// Channels for lifecycle and messaging
	register   chan *types.Client
	unregister chan *types.Client
	broadcast  chan types.Event
}

// NewHub creates and returns a new WebSocket Hub instance.
func NewHub() *Hub {
	return &Hub{
		channels:   make(map[uuid.UUID]map[*types.Client]bool),
		register:   make(chan *types.Client),
		unregister: make(chan *types.Client),
		broadcast:  make(chan types.Event),
	}
}

// Run starts the Hub's main loop to handle registration, unregistration, and broadcasting.
func (h *Hub) Run() {
	for {
		select {

		// Register a client to a channel
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.channels[client.ChannelID]; !ok {
				h.channels[client.ChannelID] = make(map[*types.Client]bool)
			}
			h.channels[client.ChannelID][client] = true
			h.mu.Unlock()

		// Unregister a client and clean up
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.channels[client.ChannelID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.channels, client.ChannelID)
					}
				}
			}
			h.mu.Unlock()

		// Broadcast an event to all clients in the channel
		case event := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.channels[event.ChannelID]; ok {
				for client := range clients {
					select {
					case client.Send <- event.Data:
					default:
						// Client is unresponsive
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage sends an event to all connected clients in the event's channel.
func (h *Hub) BroadcastMessage(event types.Event) {
	h.broadcast <- event
}

// RegisterClient adds a new client to the hub.
func (h *Hub) RegisterClient(client *types.Client) {
	h.register <- client
}

// UnregisterClient removes a client from the hub.
func (h *Hub) UnregisterClient(client *types.Client) {
	h.unregister <- client
}
