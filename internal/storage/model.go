package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Step      int          `json:"step"`
	Repo      string       `json:"repo"`
	Started   time.Time    `json:"started"`
	Submitted sql.NullTime `json:"submitted,omitempty"`
	Timelimit int          `json:"timelimit"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	IsAdmin   bool      `json:"is_admin"`
}

type HydratedSession struct {
	ID        uuid.UUID
	User      User
	Step      int
	Repo      string
	Started   time.Time
	Submitted sql.NullTime
	Timelimit int
}
