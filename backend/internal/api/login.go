package api

import (
	"encoding/json"
	"mana/internal/auth"
	"mana/internal/models"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) isValid() bool {
	hasUsername := strings.TrimSpace(req.Username) != ""
	hasEmail := strings.TrimSpace(req.Email) != ""

	if hasUsername == hasEmail {
		return false
	}

	return strings.TrimSpace(req.Password) != ""
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !req.isValid() {
		http.Error(w, "Provide either username or email, not both", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	ctx := r.Context()
	var user *models.User
	var err error

	if req.Username != "" {
		if !usernameRegex.MatchString(req.Username) {
			http.Error(w, "Invalid username format", http.StatusBadRequest)
			return
		}

		user, err = api.Store.Users.GetUserByUsername(ctx, req.Username)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Username not found", http.StatusNotFound)
			return
		}

	} else {
		if !emailRegex.MatchString(req.Email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		user, err = api.Store.Users.GetUserByEmail(ctx, req.Email)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Email not found", http.StatusNotFound)
			return
		}
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"user":  user.ToPublicUser(),
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
