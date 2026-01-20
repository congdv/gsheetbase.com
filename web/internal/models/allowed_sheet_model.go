package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AllowedSheet struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	UserID         uuid.UUID      `db:"user_id" json:"user_id"`
	SheetID        string         `db:"sheet_id" json:"sheet_id"`
	SheetName      *string        `db:"sheet_name" json:"sheet_name,omitempty"`
	Description    *string        `db:"description" json:"description,omitempty"`
	AllowWrite     bool           `db:"allow_write" json:"allow_write"`
	AllowedMethods pq.StringArray `db:"allowed_methods" json:"allowed_methods"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}
