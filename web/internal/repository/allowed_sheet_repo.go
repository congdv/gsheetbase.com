package repository

import (
	"context"
	"time"

	"gsheetbase/web/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AllowedSheetRepo interface {
	Register(ctx context.Context, userID uuid.UUID, sheetID, sheetName, description string) (models.AllowedSheet, error)
	IsAllowed(ctx context.Context, userID uuid.UUID, sheetID string) (bool, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AllowedSheet, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.AllowedSheet, error)
	Delete(ctx context.Context, userID uuid.UUID, sheetID string) error
	UpdateAuth(ctx context.Context, sheetID uuid.UUID, authType string, bearerToken, basicUsername, basicPasswordHash *string) error
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

func (r *allowedSheetRepo) Delete(ctx context.Context, userID uuid.UUID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM allowed_sheets WHERE user_id = $1 AND sheet_id = $2
	`, userID, sheetID)
	return err
}

// UpdateAuth updates the authentication type and credentials for a sheet
func (r *allowedSheetRepo) UpdateAuth(ctx context.Context, sheetID uuid.UUID, authType string, bearerToken, basicUsername, basicPasswordHash *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET auth_type = $1,
			auth_bearer_token = NULLIF($2, ''),
			auth_basic_username = NULLIF($3, ''),
			auth_basic_password_hash = NULLIF($4, ''),
			updated_at = NOW()
		WHERE id = $5
	`, authType, bearerToken, basicUsername, basicPasswordHash, sheetID)
	return err
}
