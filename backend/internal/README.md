internal/
├── api/           // HTTP route handlers
├── websocket/     // WebSocket connection & event handling
├── services/      // Core business/domain logic
│   └── message.go
├── models/        // Entity definitions
├── store/         // DB interfaces and SQL adapters
├── types/         // Shared low-level types (DTOs, enums, etc.)
├── permissions/   // Bitflag perms, resolution logic
├── auth/          // Auth helpers (tokens, sessions)
├── middleware/    // HTTP middleware (logging, auth, etc.)
├── db/            // DB init & connection logic
