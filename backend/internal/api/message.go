package api

import (
	"encoding/json"
	"mana/internal/middleware"
	"mana/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MessageContent struct {
	Content string `json:"content"`
}

func (api *API) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get our channel id
	channelIDParam := chi.URLParam(r, "id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// make sure they are authorized
	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input MessageContent

	// get our input message
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Invalid message content", http.StatusBadRequest)
		return
	}

	// insert message
	msg := models.NewMessage(channelID, userID, input.Content)
	if err := api.Store.Messages.InsertMessage(ctx, msg); err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (api *API) GetMessagesByChannel(w http.ResponseWriter, r *http.Request) {
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

// func (api *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {}

//func (api *API) EditMessage(w http.ResponseWriter, r *http.Request) {}
