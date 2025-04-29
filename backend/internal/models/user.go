package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`              // primary key, unique, UUIDv4
	Username       string    `json:"username"`        // unique
	Email          string    `json:"email"`           // unique
	Password       string    `json:"-"`               // hashed, not over API
	ActivityStatus string    `json:"activity_status"` // "online", "offline", "away", etc.
	AccountStatus  string    `json:"account_status"`  // "active", "suspended", "banned", "ip-banned"
	CreatedAt      string    `json:"created_at"`      // ISO timestamp
}

func NewUser(username string, email string, hashedPassword string) *User {
	return &User{
		ID:             uuid.New(),
		Username:       username,
		Email:          email,
		Password:       hashedPassword,
		ActivityStatus: "offline",
		AccountStatus:  "active",
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
}

func (user *User) SetActivityStatus(status string) {
	user.ActivityStatus = status
}

func (user *User) SetAccountStatus(status string) {
	user.AccountStatus = status
}
