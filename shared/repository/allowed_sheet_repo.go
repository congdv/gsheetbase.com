package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"gsheetbase/shared/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AllowedSheetRepo interface {
	Register(ctx context.Context, userID uuid.UUID, sheetID, sheetName, description string) (models.AllowedSheet, error)
	IsAllowed(ctx context.Context, userID uuid.UUID, sheetID string) (bool, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AllowedSheet, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.AllowedSheet, error)
	FindByAPIKey(ctx context.Context, apiKey string) (models.AllowedSheet, error)
	Delete(ctx context.Context, userID uuid.UUID, sheetID string) error
	Publish(ctx context.Context, sheetID uuid.UUID, defaultRange string, useFirstRowAsHeader bool) (string, error)
	Unpublish(ctx context.Context, sheetID uuid.UUID) error
	UpdateWriteSettings(ctx context.Context, sheetID uuid.UUID, allowWrite bool) error
}

type allowedSheetRepo struct {
	db *sqlx.DB
}

func NewAllowedSheetRepo(db *sqlx.DB) AllowedSheetRepo {
	return &allowedSheetRepo{db: db}
}

func (r *allowedSheetRepo) Register(ctx context.Context, userID uuid.UUID, sheetID, sheetName, description string) (models.AllowedSheet, error) {
	now := time.Now()
	sheet := models.AllowedSheet{
		ID:        uuid.New(),
		UserID:    userID,
		SheetID:   sheetID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if sheetName != "" {
		sheet.SheetName = &sheetName
	}
	if description != "" {
		sheet.Description = &description
	}

	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO allowed_sheets (id, user_id, sheet_id, sheet_name, description, created_at, updated_at)
		VALUES (:id, :user_id, :sheet_id, :sheet_name, :description, :created_at, :updated_at)
		ON CONFLICT (user_id, sheet_id) DO UPDATE 
		SET sheet_name = EXCLUDED.sheet_name, 
		    description = EXCLUDED.description,
		    updated_at = EXCLUDED.updated_at
		RETURNING *
	`, &sheet)

	return sheet, err
}

func (r *allowedSheetRepo) IsAllowed(ctx context.Context, userID uuid.UUID, sheetID string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM allowed_sheets WHERE user_id = $1 AND sheet_id = $2
	`, userID, sheetID)
	return count > 0, err
}

func (r *allowedSheetRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AllowedSheet, error) {
	var sheets []models.AllowedSheet
	err := r.db.SelectContext(ctx, &sheets, `
		SELECT * FROM allowed_sheets WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	return sheets, err
}

func (r *allowedSheetRepo) FindByID(ctx context.Context, id uuid.UUID) (models.AllowedSheet, error) {
	var sheet models.AllowedSheet
	err := r.db.GetContext(ctx, &sheet, `
		SELECT * FROM allowed_sheets WHERE id = $1
	`, id)
	return sheet, err
}

func (r *allowedSheetRepo) FindByAPIKey(ctx context.Context, apiKey string) (models.AllowedSheet, error) {
	var sheet models.AllowedSheet
	err := r.db.GetContext(ctx, &sheet, `
		SELECT * FROM allowed_sheets WHERE api_key = $1 AND is_public = true
	`, apiKey)
	return sheet, err
}

func (r *allowedSheetRepo) Delete(ctx context.Context, userID uuid.UUID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM allowed_sheets WHERE user_id = $1 AND sheet_id = $2
	`, userID, sheetID)
	return err
}

func (r *allowedSheetRepo) Publish(ctx context.Context, sheetID uuid.UUID, defaultRange string, useFirstRowAsHeader bool) (string, error) {
	// Generate API key
	apiKey := generateAPIKey()

	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET api_key = $1, 
		    is_public = true,
		    default_range = $2,
		    use_first_row_as_header = $3,
		    updated_at = NOW()
		WHERE id = $4
	`, apiKey, defaultRange, useFirstRowAsHeader, sheetID)

	return apiKey, err
}

func (r *allowedSheetRepo) Unpublish(ctx context.Context, sheetID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET api_key = NULL, 
		    is_public = false,
		    updated_at = NOW()
		WHERE id = $1
	`, sheetID)
	return err
}

func (r *allowedSheetRepo) UpdateWriteSettings(ctx context.Context, sheetID uuid.UUID, allowWrite bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET allow_write = $1,
		    updated_at = NOW()
		WHERE id = $2
	`, allowWrite, sheetID)
	return err
}

func generateAPIKey() string {
	b := make([]byte, 24)
	rand.Read(b)
	return "gsheet_" + base64.URLEncoding.EncodeToString(b)[:32]
}
