package websocket

import (
	"bytes"
	"log"
	"time"

	"mana/internal/types"
	"mana/internal/websocket/events"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type ClientImpl struct {
	Client     types.Client
	Connection *websocket.Conn
}

// pump messages incoming from socket to the hub (inbound messages)
func (client *ClientImpl) readPump() {

	// Deinit hub and close connect on end
	defer func() {
		client.Client.Hub.UnregisterClient(&client.Client)
		client.Connection.Close()
	}()

	// Setup connection
	client.Connection.SetReadLimit(maxMessageSize)
	client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	client.Connection.SetPongHandler(func(appdata string) error {
		client.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		// read inbound messages
		_, message, err := client.Connection.ReadMessage()

		// unrecoverable error
		if err != nil {
			// if error closes connection, log error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		go events.HandleEvent(&client.Client, message)
	}
}

// pump message from our socket to hub (outbound message)
func (client *ClientImpl) writePump() {
	ticker := time.NewTicker(pingPeriod)

	// deinit ticker and connection
	defer func() {
		ticker.Stop()
		client.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-client.Client.Send:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed channel
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writer.Write(message)

			// drain queued messages (avoid blocks)
			n := len(client.Client.Send)
			for i := 0; i < n; i++ {
				_, _ = writer.Write([]byte("\n"))
				_, _ = writer.Write(<-client.Client.Send)
			}

			// if our writer closed on err
			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
