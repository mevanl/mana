package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) GetGuildChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildIDParam := chi.URLParam(r, "id")
	guildID, err := uuid.Parse(guildIDParam)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	channels, err := api.Store.GuildChannels.GetChannelsForGuild(ctx, guildID)
	if err != nil {
		http.Error(w, "Failed to fetch channels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}

//func (api *API) EditChannelName(w http.ResponseWriter, r *http.Request) {}

//func (api *API) EditChannelTopic(w http.ResponseWriter, r *http.Request) {}
