package api

import (
	"encoding/json"
	"mana/internal/auth"
	"mana/internal/models"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{5,32}$`)
var passwordRegex = regexp.MustCompile(`^[A-Za-z0-9?.=*!]{8,64}$`)
var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// decode request body and place into our req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// clean whitespace
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// validation
	if !usernameRegex.MatchString(req.Username) {
		http.Error(w, "Invalid username. Must be: 5-32 characters, alphanumeric or underscore only", http.StatusBadRequest)
		return
	}

	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "Invalid email format.", http.StatusBadRequest)
		return
	}

	if !checkValidPassword(req.Password) {
		http.Error(w, "Password must be 8-64 characters, including: 1 uppercase, 1 lowercase, 1 digit.", http.StatusBadRequest)
		return
	}

	// Get request context
	ctx := r.Context()

	// Check store existance for email and username
	if exists, _ := api.Store.Users.CheckUserExistsByEmail(ctx, req.Email); exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	if exists, _ := api.Store.Users.CheckUserExistsByUsername(ctx, req.Username); exists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal error hashing password", http.StatusInternalServerError)
		return
	}

	// Make user and insert into db
	user := models.NewUser(req.Username, req.Email, hashedPassword)

	err = api.Store.Users.Create(ctx, user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// generate user JWT
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

func checkValidPassword(password string) bool {
	if match := passwordRegex.MatchString(password); !match {
		return false
	}

	var hasUpper bool = false
	var hasLower bool = false
	var hasDigit bool = false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
