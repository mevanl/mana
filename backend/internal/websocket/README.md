internal/
├── websocket/
│   ├── websocket.go       // HTTP handler: upgrades request, auths user, registers client
│   ├── hub.go             // Global Hub: manages connected clients, broadcasts, etc.
│   ├── client.go          // Client struct: per-connection socket, send/recv routines
│   ├── router.go          // Dispatch incoming events to handlers
│   └── handlers/          // One file per event type
│       ├── message.go     // e.g., handleSendMessage
│       ├── typing.go      // e.g., handleStartTyping
│       └── ...
