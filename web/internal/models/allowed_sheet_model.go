package models

import (
	"time"

	"github.com/google/uuid"
)

type AllowedSheet struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	SheetID     string    `db:"sheet_id" json:"sheet_id"`
	SheetName   *string   `db:"sheet_name" json:"sheet_name,omitempty"`
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
