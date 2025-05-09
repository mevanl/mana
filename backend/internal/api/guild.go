package api

import (
	"encoding/json"
	"mana/internal/middleware"
	"mana/internal/models"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateGuildRequest struct {
	Name string `json:"name"`
}

type GetGuildByIDRequest struct {
	ID string `json:"id"`
}

func (api *API) CreateGuild(w http.ResponseWriter, r *http.Request) {
	var req CreateGuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// input vaidation
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) < 2 || len(req.Name) > 100 {
		http.Error(w, "Guild name must be between 2 and 100 characters.", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(uuid.UUID)

	// Create guild and owner membership
	guild := models.NewGuild(req.Name, userID)

	// insert into db
	if err := api.Store.Guilds.InsertGuild(ctx, guild); err != nil {
		http.Error(w, "Failed to create guild", http.StatusInternalServerError)
		return
	}

	// TODO: handle everyone/default role ? handle here or somewhere else

	resp := map[string]interface{}{
		"guild": guild,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	guild, err := api.Store.Guilds.GetGuildByID(ctx, guildID)
	if err != nil {
		http.Error(w, "Could not find guild.", http.StatusNotFound)
		return
	}

	// success
	resp := map[string]interface{}{"guild": guild}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (api *API) DeleteGuild(w http.ResponseWriter, r *http.Request) {
	// grab THAT id
	guildIDStr := chi.URLParam(r, "id")
	guildID, err := uuid.Parse(guildIDStr)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(uuid.UUID)

	guild, err := api.Store.Guilds.GetGuildByID(ctx, guildID)
	if err != nil {
		http.Error(w, "Guild not found", http.StatusNotFound)
		return
	}

	if guild.OwnerID != userID {
		http.Error(w, "You do not have permission to delete this guild", http.StatusForbidden)
		return
	}

	if err := api.Store.Guilds.DeleteGuild(ctx, guildID); err != nil {
		http.Error(w, "Failed to delete guild", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) GetUserGuilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	guilds, err := api.Store.Guilds.GetGuildsForUserID(ctx, userID)
	if err != nil {
		http.Error(w, "Failed to fetch guilds", http.StatusInternalServerError)
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

	// Get guild by invite code
	guild, err := api.Store.Guilds.GetGuildByInviteCode(ctx, inviteCode)
	if err != nil {
		http.Error(w, "Invalid or expired invite code", http.StatusNotFound)
		return
	}

	// Check if user is member of guild already
	exists, err := api.Store.Guilds.CheckUserMemberOfGuild(ctx, guild.ID, userID)
	if err != nil {
		http.Error(w, "Failed to check membership", http.StatusInternalServerError)
		return
	}

	// already a member
	if exists {
		http.Error(w, "Already a member of this guild", http.StatusConflict)
		return
	}

	// Add user to guild
	member := models.NewGuildMember(guild.ID, userID)
	err = api.Store.Guilds.AddUserToGuild(ctx, member)
	if err != nil {
		http.Error(w, "Failed to join guild", http.StatusInternalServerError)
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
