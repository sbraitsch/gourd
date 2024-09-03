package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uuid.UUID    `json:"id"`
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Step      int          `json:"step"`
	Started   time.Time    `json:"started"`
	Submitted sql.NullTime `json:"submitted,omitempty"`
	Timelimit int          `json:"timelimit"`
}
