package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) GetChannelMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get channel 
	channelIDParam := chi.URLParam(r, "id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Get limit from query string 
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// grab the messages from our guild 
	messages, err := api.Store.Messages.GetMessagesByChannel(ctx, channelID, limit)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// send the messages 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

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
