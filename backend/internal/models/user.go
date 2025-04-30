package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`                        // primary key, unique, UUIDv4
	Username       string    `json:"username"`                  // unique
	Email          string    `json:"email"`                     // unique
	Password       string    `json:"-"`                         // hashed, not over API
	ActivityStatus string    `json:"activity_status,omitempty"` // "online", "offline", "away", etc.
	AccountStatus  string    `json:"account_status,omitempty"`  // "active", "suspended", "banned"
	CreatedAt      time.Time `json:"created_at"`                // ISO timestamp
}

func NewUser(username string, email string, hashedPassword string) *User {
	return &User{
		ID:             uuid.New(),
		Username:       username,
		Email:          email,
		Password:       hashedPassword,
		ActivityStatus: "offline",
		AccountStatus:  "active",
		CreatedAt:      time.Now().UTC(),
	}
}

func (user *User) SetActivityStatus(status string) {
	user.ActivityStatus = status
}

func (user *User) SetAccountStatus(status string) {
	user.AccountStatus = status
}

// Sanitized version of a user for public view
type PublicUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	ActivityStatus string    `json:"activity_status,omitempty"`
}

// Convert User to PublicUser
func (user *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:             user.ID,
		Username:       user.Username,
		ActivityStatus: user.ActivityStatus,
	}
}

// Convert an array of unsanitized users to public users
func ToPublicUserSlice(users []*User) []*PublicUser {
	publicUsers := make([]*PublicUser, len(users))

	for i, user := range users {
		publicUsers[i] = user.ToPublicUser()
	}

	return publicUsers
}
