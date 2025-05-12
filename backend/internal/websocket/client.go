package websocket

import (
	"bytes"
	"log"
	"time"

	"mana/internal/types"
	"mana/internal/websocket/router"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// ClientImpl handles the WebSocket connection and communication.
type ClientImpl struct {
	Client     types.Client
	Connection *websocket.Conn
}

// readPump listens for incoming WebSocket messages and dispatches them to the router.
func (c *ClientImpl) readPump() {
	defer func() {
		c.Client.Hub.UnregisterClient(&c.Client)
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket: unexpected close: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		go router.HandleEvent(&c.Client, message)
	}
}

// writePump sends messages from the client.Send channel to the WebSocket connection.
func (c *ClientImpl) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Client.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				_ = c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("websocket: writer error: %v", err)
				return
			}

			_, _ = writer.Write(message)

			// Drain queued messages (non-blocking)
			n := len(c.Client.Send)
			for i := 0; i < n; i++ {
				_, _ = writer.Write([]byte("\n"))
				_, _ = writer.Write(<-c.Client.Send)
			}

			if err := writer.Close(); err != nil {
				log.Printf("websocket: writer close error: %v", err)
				return
			}

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("websocket: ping error: %v", err)
				return
			}
		}
	}
}
