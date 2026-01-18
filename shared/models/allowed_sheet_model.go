package models

import (
	"time"

	"github.com/google/uuid"
)

type AllowedSheet struct {
	ID                     uuid.UUID `db:"id" json:"id"`
	UserID                 uuid.UUID `db:"user_id" json:"user_id"`
	SheetID                string    `db:"sheet_id" json:"sheet_id"`
	SheetName              *string   `db:"sheet_name" json:"sheet_name,omitempty"`
	Description            *string   `db:"description" json:"description,omitempty"`
	APIKey                 *string   `db:"api_key" json:"api_key,omitempty"`
	IsPublic               bool      `db:"is_public" json:"is_public"`
	DefaultRange           *string   `db:"default_range" json:"default_range,omitempty"`
	UseFirstRowAsHeader    bool      `db:"use_first_row_as_header" json:"use_first_row_as_header"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}
