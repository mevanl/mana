package handler

import (
	"mana/internal/store"
)

type Handler struct {
	store *store.Store
	// In future: GuildStore, UserStore, etc.
}
