package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AllowedSheet struct {
	ID                    uuid.UUID      `db:"id" json:"id"`
	UserID                uuid.UUID      `db:"user_id" json:"user_id"`
	SheetID               string         `db:"sheet_id" json:"sheet_id"`
	SheetName             *string        `db:"sheet_name" json:"sheet_name,omitempty"`
	Description           *string        `db:"description" json:"description,omitempty"`
	APIKey                *string        `db:"api_key" json:"api_key,omitempty"`
	IsPublic              bool           `db:"is_public" json:"is_public"`
	DefaultRange          *string        `db:"default_range" json:"default_range,omitempty"`
	UseFirstRowAsHeader   bool           `db:"use_first_row_as_header" json:"use_first_row_as_header"`
	AllowWrite            bool           `db:"allow_write" json:"allow_write"`
	AllowedMethods        pq.StringArray `db:"allowed_methods" json:"allowed_methods"`
	AuthType              string         `db:"auth_type" json:"auth_type"`
	AuthBearerToken       *string        `db:"auth_bearer_token" json:"auth_bearer_token,omitempty"`
	AuthBasicUsername     *string        `db:"auth_basic_username" json:"auth_basic_username,omitempty"`
	AuthBasicPasswordHash *string        `db:"auth_basic_password_hash" json:"-"`
	CreatedAt             time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at" json:"updated_at"`
}
