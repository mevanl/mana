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
		// will need implentation, keep open for dev testing for now
		return true
	},
}

func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// get the channel id user is in
	channelIDStr := chi.URLParam(r, "channel_id")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel_id", http.StatusBadRequest)
		return
	}

	// authenticate user
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// upgrade the socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	log.Printf("WebSocket connected: user=%s channel=%s\n", userID, channelID)

	// create our client
	client := &ClientImpl{
		Client: types.Client{
			Ctx:       r.Context(),
			Hub:       hub,
			Send:      make(chan []byte, 256),
			UserID:    userID,
			ChannelID: channelID,
		},
		Connection: conn,
	}

	conn.SetCloseHandler(func(code int, text string) error {
		hub.unregister <- &client.Client
		return nil
	})

	// register them
	hub.register <- &client.Client

	// start goroutines
	go client.writePump()
	go client.readPump()
}
