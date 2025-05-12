package types

import "mana/internal/models"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User  models.PublicUser `json:"user"`
	Token string            `json:"token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  models.PublicUser `json:"user"`
	Token string            `json:"token"`
}

// Next: Suggested DTOs to Add Soon

//     CreateGuildRequest

//     CreateChannelRequest

//     SendMessageRequest

//     UpdateUserStatusRequest
