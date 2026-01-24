package models

import (
	"time"

	"github.com/google/uuid"
)

// ApiUsageDaily represents daily API usage statistics
type ApiUsageDaily struct {
	ID           uuid.UUID `db:"id"`
	ApiKey       string    `db:"api_key"`
	UserID       uuid.UUID `db:"user_id"`
	SheetID      uuid.UUID `db:"sheet_id"`
	RequestDate  time.Time `db:"request_date"`
	Method       string    `db:"method"`
	RequestCount int       `db:"request_count"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
