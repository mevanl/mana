package api

import (
	manaerror "mana/internal/errors"
	"net/http"
)

func RespondGuildError(w http.ResponseWriter, err error) {
	switch {
	// Guild Errors
	case manaerror.Is(err, manaerror.ErrGuildNotFound):
		http.Error(w, "Guild not found", http.StatusNotFound)
	case manaerror.Is(err, manaerror.ErrNotGuildOwner):
		http.Error(w, "You do not have permission to delete this guild", http.StatusForbidden)
	case manaerror.Is(err, manaerror.ErrGuildDeleteFailed):
		http.Error(w, "Failed to delete guild", http.StatusInternalServerError)
	case manaerror.Is(err, manaerror.ErrBadGuildName):
		http.Error(w, "invalid guild name", http.StatusBadRequest)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func RespondMessageError(w http.ResponseWriter, err error) {
	switch {
	case manaerror.Is(err, manaerror.ErrMessageSendFailed):
		http.Error(w, "message could not send", http.StatusInternalServerError)
	case manaerror.Is(err, manaerror.ErrMessageEmpty):
		http.Error(w, "empty message", http.StatusBadRequest)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func RespondChannelError(w http.ResponseWriter, err error) {
	switch {
	case manaerror.Is(err, manaerror.ErrChannelFetchFailed):
		http.Error(w, "Could not fetch channels", http.StatusInternalServerError)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
