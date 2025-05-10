package api

import (
	"encoding/json"
	"mana/internal/middleware"
	"mana/internal/models"
	"net/http"

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

// func (api *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {}

//func (api *API) EditMessage(w http.ResponseWriter, r *http.Request) {}
