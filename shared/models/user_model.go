package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	Email              string         `db:"email" json:"email"`
	Provider           string         `db:"provider" json:"-"`
	ProviderID         string         `db:"provider_id" json:"-"`
	GoogleAccessToken  *string        `db:"google_access_token" json:"google_access_token,omitempty"`
	GoogleRefreshToken *string        `db:"google_refresh_token" json:"-"`
	GoogleTokenExpiry  *time.Time     `db:"google_token_expiry" json:"google_token_expiry,omitempty"`
	GoogleScopes       pq.StringArray `db:"google_scopes" json:"google_scopes,omitempty"`
	RateLimitPerMinute int            `db:"rate_limit_per_minute" json:"rate_limit_per_minute"`
	RateLimitBurst     int            `db:"rate_limit_burst" json:"rate_limit_burst"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at" json:"updated_at"`
}
