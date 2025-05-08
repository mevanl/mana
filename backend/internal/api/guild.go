package api

import (
	"encoding/json"
	"mana/internal/models"
	"net/http"
	"strings"

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
	userID := ctx.Value("userID").(uuid.UUID)

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

func GetGuildByID(w http.ResponseWriter, r *http.Request) {

}
