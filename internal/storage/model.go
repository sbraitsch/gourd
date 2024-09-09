package storage

import (
	"database/sql"
	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	CurrentStep int          `json:"step"`
	MaxProgress int          `json:"max_progress"`
	Repo        string       `json:"repo"`
	Started     sql.NullTime `json:"started, omitempty"`
	Submitted   sql.NullTime `json:"submitted,omitempty"`
	Timelimit   int          `json:"timelimit"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	IsAdmin   bool      `json:"is_admin"`
}

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
