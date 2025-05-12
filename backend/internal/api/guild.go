package api

import (
	"encoding/json"
	"mana/internal/middleware"
	"mana/internal/services"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateGuildRequest struct {
	Name string `json:"name"`
}

func (api *API) CreateGuild(w http.ResponseWriter, r *http.Request) {
	var req CreateGuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(uuid.UUID)

	guild, err := services.CreateGuild(ctx, api.Store, req.Name, userID)
	if err != nil {
		RespondGuildError(w, err)
		return
	}

	resp := map[string]interface{}{
		"guild": guild,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (api *API) DeleteGuild(w http.ResponseWriter, r *http.Request) {
	// grab id
	guildIDStr := chi.URLParam(r, "id")
	guildID, err := uuid.Parse(guildIDStr)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(uuid.UUID)

	err = services.DeleteGuild(ctx, api.Store, guildID, userID)
	if err != nil {
		RespondGuildError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) GetGuildByID(w http.ResponseWriter, r *http.Request) {
	// grab that id
	guildIDStr := chi.URLParam(r, "id")
	if guildIDStr == "" {
		http.Error(w, "Missing guild id", http.StatusBadRequest)
		return
	}

	// convert str -> uuid
	guildID, err := uuid.Parse(strings.TrimSpace(guildIDStr))
	if err != nil {
		http.Error(w, "Invalid guild id", http.StatusBadRequest)
		return
	}

	// find that guild
	ctx := r.Context()
	guild, err := services.GetGuild(ctx, api.Store, guildID)
	if err != nil {
		RespondGuildError(w, err)
	}

	// success
	resp := map[string]interface{}{"guild": guild}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (api *API) GetUserGuilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	guilds, err := services.GetUserGuilds(ctx, api.Store, userID)
	if err != nil {
		RespondGuildError(w, err)
		return
	}

	resp := map[string]interface{}{
		"guilds": guilds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (api *API) JoinGuildByInvite(w http.ResponseWriter, r *http.Request) {
	inviteCode := chi.URLParam(r, "code")

	// empty code
	if inviteCode == "" {
		http.Error(w, "Missing invite code", http.StatusBadRequest)
		return
	}

	// grab user id from ctx
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	guild, err := services.JoinGuildByInvite(ctx, api.Store, userID, inviteCode)
	if err != nil {
		RespondGuildError(w, err)
		return
	}

	// do @everyone here ?

	// send response
	resp := map[string]interface{}{
		"guild": guild,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
