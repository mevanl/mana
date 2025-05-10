package websocket

import (
	"log"
	"mana/internal/auth"
	"mana/internal/types"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// we are allowing all origins for now
		return true
	},
}

func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// extra channel ID from query
	channelIDStr := chi.URLParam(r, "channel_id")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel_id", http.StatusBadRequest)
		return
	}

	// get user id
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// upgrade connection
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	log.Printf("WebSocket connected: user=%s channel=%s\n", userID, channelID)

	client := &ClientImpl{
		Client: types.Client{
			Hub:       hub,
			Send:      make(chan []byte, 256),
			UserID:    userID,
			ChannelID: channelID,
		},
		Connection: connection,
	}

	connection.SetCloseHandler(func(code int, text string) error {
		hub.Unregister <- &client.Client
		return nil
	})

	hub.Register <- &client.Client

	// Start pumps
	go client.writePump()
	go client.readPump()
}
