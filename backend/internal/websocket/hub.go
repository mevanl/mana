package websocket

import (
	"mana/internal/types"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	// read write lock
	mutex sync.RWMutex

	// maps channel id to all clients connected to that channel
	Channels map[uuid.UUID]map[*types.Client]bool

	// Channels for events
	Register   chan *types.Client
	Unregister chan *types.Client
	Broadcast  chan types.Event
}

func NewHub() *Hub {
	return &Hub{
		Channels:   make(map[uuid.UUID]map[*types.Client]bool),
		Register:   make(chan *types.Client),
		Unregister: make(chan *types.Client),
		Broadcast:  make(chan types.Event),
	}
}

// run is the event loop that will listen for all hub actions
func (hub *Hub) Run() {
	for {

		select {

		// Register a client to a channel
		case client := <-hub.Register:
			hub.mutex.Lock()

			if _, ok := hub.Channels[client.ChannelID]; !ok {

				// initialize that given channel entry and map it to the client
				hub.Channels[client.ChannelID] = make(map[*types.Client]bool)
			}

			hub.Channels[client.ChannelID][client] = true

			hub.mutex.Unlock()

		// remove client from the set of clients for this channel
		case client := <-hub.Unregister:
			hub.mutex.Lock()

			// if we have clients for that channel
			if clients, ok := hub.Channels[client.ChannelID]; ok {

				// if this specific client exists in all the clients for that channel
				if _, exists := clients[client]; exists {

					// unregister them from clients
					delete(clients, client)
					close(client.Send)

					// if we have 0 clients left, remove this channel from hub
					if len(clients) == 0 {
						delete(hub.Channels, client.ChannelID)
					}
				}
			}

			hub.mutex.Unlock()

		// Deliver an event to all clients connected to a channelID
		case event := <-hub.Broadcast:
			hub.mutex.Lock()

			// if we have clients in the given event's channel
			if clients, ok := hub.Channels[event.ChannelID]; ok {

				// send event data to every connected client for that channel
				for client := range clients {
					select {
					case client.Send <- event.Data:
					default:
						// client is unresponsive
						close(client.Send)
						delete(clients, client)
					}
				}
			}

			hub.mutex.Unlock()
		}

	}
}

func (h *Hub) BroadcastMessage(event types.Event) {
	h.Broadcast <- event
}

func (h *Hub) UnregisterClient(client *types.Client) {
	h.Unregister <- client
}

func (h *Hub) RegisterClient(client *types.Client) {
	h.Register <- client
}
