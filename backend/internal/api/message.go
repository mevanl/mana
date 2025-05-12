package api

import (
	"encoding/json"
	"mana/internal/middleware"
	"mana/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MessageContent struct {
	Content string `json:"content"`
}

func (api *API) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract channel ID from URL
	channelIDParam := chi.URLParam(r, "id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Ensure user is authenticated
	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var input MessageContent
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Content == "" {
		http.Error(w, "Invalid message content", http.StatusBadRequest)
		return
	}

	// Call the centralized service logic
	msg, err := services.SendMessage(ctx, api.Store, userID, channelID, input.Content)
	if err != nil {
		RespondMessageError(w, err)
	}

	// Return the message as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (api *API) GetMessagesByChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelIDParam := chi.URLParam(r, "id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var before *time.Time
	if b := r.URL.Query().Get("before"); b != "" {
		t, err := time.Parse(time.RFC3339, b)
		if err != nil {
			http.Error(w, "Invalid before timestamp", http.StatusBadRequest)
			return
		}
		before = &t
	}

	messages, err := services.GetChannelMessages(ctx, api.Store, channelID, limit, before)
	if err != nil {
		RespondMessageError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// func (api *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {}

//func (api *API) EditMessage(w http.ResponseWriter, r *http.Request) {}
