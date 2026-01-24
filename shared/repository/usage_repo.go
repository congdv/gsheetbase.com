package repository

import (
	"context"
	"time"

	"gsheetbase/shared/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UsageRepo defines the interface for usage tracking operations
type UsageRepo interface {
	IncrementDailyUsage(ctx context.Context, apiKey string, userID, sheetID uuid.UUID, date time.Time, method string) error
	GetDailyUsageBySheet(ctx context.Context, sheetID uuid.UUID, startDate, endDate time.Time) ([]models.ApiUsageDaily, error)
	GetDailyUsageByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.ApiUsageDaily, error)
	GetDailyUsageByAPIKey(ctx context.Context, apiKey string, startDate, endDate time.Time) ([]models.ApiUsageDaily, error)

	// Quota checking methods
	GetTodayUsageCount(ctx context.Context, userID uuid.UUID, method string) (int, error)
	GetMonthlyUsageCount(ctx context.Context, userID uuid.UUID, method string) (int, error)
}

type usageRepo struct {
	db *sqlx.DB
}

// NewUsageRepo creates a new usage repository
func NewUsageRepo(db *sqlx.DB) UsageRepo {
	return &usageRepo{db: db}
}

// IncrementDailyUsage atomically increments the usage counter
func (r *usageRepo) IncrementDailyUsage(ctx context.Context, apiKey string, userID, sheetID uuid.UUID, date time.Time, method string) error {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	query := `
		INSERT INTO api_usage_daily (api_key, user_id, sheet_id, request_date, method, request_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 1, NOW(), NOW())
		ON CONFLICT (api_key, request_date, method)
		DO UPDATE SET 
			request_count = api_usage_daily.request_count + 1,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, apiKey, userID, sheetID, dateOnly, method)
	return err
}

// GetDailyUsageBySheet retrieves usage stats for a specific sheet
func (r *usageRepo) GetDailyUsageBySheet(ctx context.Context, sheetID uuid.UUID, startDate, endDate time.Time) ([]models.ApiUsageDaily, error) {
	query := `
		SELECT id, api_key, user_id, sheet_id, request_date, method, request_count, created_at, updated_at
		FROM api_usage_daily
		WHERE sheet_id = $1 AND request_date >= $2 AND request_date <= $3
		ORDER BY request_date DESC, method
	`

	rows, err := r.db.QueryContext(ctx, query, sheetID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ApiUsageDaily
	for rows.Next() {
		var usage models.ApiUsageDaily
		err := rows.Scan(&usage.ID, &usage.ApiKey, &usage.UserID, &usage.SheetID, &usage.RequestDate, &usage.Method, &usage.RequestCount, &usage.CreatedAt, &usage.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, usage)
	}

	return results, rows.Err()
}

// GetDailyUsageByUser retrieves usage stats for all sheets owned by a user
func (r *usageRepo) GetDailyUsageByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.ApiUsageDaily, error) {
	query := `
		SELECT id, api_key, user_id, sheet_id, request_date, method, request_count, created_at, updated_at
		FROM api_usage_daily
		WHERE user_id = $1 AND request_date >= $2 AND request_date <= $3
		ORDER BY request_date DESC, sheet_id, method
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ApiUsageDaily
	for rows.Next() {
		var usage models.ApiUsageDaily
		err := rows.Scan(&usage.ID, &usage.ApiKey, &usage.UserID, &usage.SheetID, &usage.RequestDate, &usage.Method, &usage.RequestCount, &usage.CreatedAt, &usage.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, usage)
	}

	return results, rows.Err()
}

// GetDailyUsageByAPIKey retrieves usage stats for a specific API key
func (r *usageRepo) GetDailyUsageByAPIKey(ctx context.Context, apiKey string, startDate, endDate time.Time) ([]models.ApiUsageDaily, error) {
	query := `
		SELECT id, api_key, user_id, sheet_id, request_date, method, request_count, created_at, updated_at
		FROM api_usage_daily
		WHERE api_key = $1 AND request_date >= $2 AND request_date <= $3
		ORDER BY request_date DESC, method
	`

	rows, err := r.db.QueryContext(ctx, query, apiKey, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ApiUsageDaily
	for rows.Next() {
		var usage models.ApiUsageDaily
		err := rows.Scan(&usage.ID, &usage.ApiKey, &usage.UserID, &usage.SheetID, &usage.RequestDate, &usage.Method, &usage.RequestCount, &usage.CreatedAt, &usage.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, usage)
	}

	return results, rows.Err()
}

// GetTodayUsageCount returns the total request count for today by user and method
func (r *usageRepo) GetTodayUsageCount(ctx context.Context, userID uuid.UUID, method string) (int, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var count int
	query := `
		SELECT COALESCE(SUM(request_count), 0)
		FROM api_usage_daily
		WHERE user_id = $1 AND request_date = $2 AND method = $3
	`

	err := r.db.GetContext(ctx, &count, query, userID, today, method)
	return count, err
}

// GetMonthlyUsageCount returns the total request count for the current month by user and method
func (r *usageRepo) GetMonthlyUsageCount(ctx context.Context, userID uuid.UUID, method string) (int, error) {
	now := time.Now().UTC()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var count int
	query := `
		SELECT COALESCE(SUM(request_count), 0)
		FROM api_usage_daily
		WHERE user_id = $1 AND request_date >= $2 AND method = $3
	`

	err := r.db.GetContext(ctx, &count, query, userID, firstOfMonth, method)
	return count, err
}
