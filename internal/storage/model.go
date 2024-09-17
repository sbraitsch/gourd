package storage

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// Session holds all information about a session in its database representation.
type Session struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	CurrentStep int          `json:"step"`
	MaxProgress int          `json:"max_progress"`
	Repo        string       `json:"repo"`
	Started     sql.NullTime `json:"started,omitempty"`
	Submitted   sql.NullTime `json:"submitted,omitempty"`
	Timelimit   int          `json:"timelimit"`
}

// User  holds all information about a user in its database representation.
type User struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	IsAdmin   bool      `json:"is_admin"`
}

// GetBranchName returns the string to be used for branch creation based on the user.
func (user *User) GetBranchName() string {
	return fmt.Sprintf("%s_%s_%s", user.Firstname, user.Lastname, user.ID)
}

// HydratedSession holds all relevant user-session information, combining session and user.
type HydratedSession struct {
	ID          uuid.UUID
	User        User
	CurrentStep int
	MaxProgress int
	Repo        string
	Started     sql.NullTime
	Submitted   sql.NullTime
	Timelimit   int
}
