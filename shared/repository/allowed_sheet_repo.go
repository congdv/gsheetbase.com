package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"gsheetbase/shared/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("unauthorized: invalid credentials")

type AllowedSheetRepo interface {
	Register(ctx context.Context, userID uuid.UUID, sheetID, sheetName, description string) (models.AllowedSheet, error)
	IsAllowed(ctx context.Context, userID uuid.UUID, sheetID string) (bool, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AllowedSheet, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.AllowedSheet, error)
	FindByAPIKey(ctx context.Context, apiKey string) (models.AllowedSheet, error)
	FindByBearerToken(ctx context.Context, token string) (models.AllowedSheet, error)
	FindByBasicCredentials(ctx context.Context, username, password string) (models.AllowedSheet, error)
	Delete(ctx context.Context, userID uuid.UUID, sheetID string) error
	Publish(ctx context.Context, sheetID uuid.UUID, defaultRange string, useFirstRowAsHeader bool) (string, error)
	Unpublish(ctx context.Context, sheetID uuid.UUID) error
	UpdateWriteSettings(ctx context.Context, sheetID uuid.UUID, allowWrite bool) error
	UpdateAllowedMethods(ctx context.Context, sheetID uuid.UUID, allowedMethods []string) error
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

func (r *allowedSheetRepo) UpdateAllowedMethods(ctx context.Context, sheetID uuid.UUID, allowedMethods []string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET allowed_methods = $1,
		    updated_at = NOW()
		WHERE id = $2
	`, pq.Array(allowedMethods), sheetID)
	return err
}

// FindByBearerToken finds a sheet by bearer token (auth_type = 'bearer')
// Bearer tokens are compared directly (stored as-is in the database).
func (r *allowedSheetRepo) FindByBearerToken(ctx context.Context, token string) (models.AllowedSheet, error) {
	var sheet models.AllowedSheet
	err := r.db.GetContext(ctx, &sheet, `
		SELECT * FROM allowed_sheets 
		WHERE auth_type = 'bearer' AND auth_bearer_token = $1 AND is_public = true
	`, token)
	return sheet, err
}

// FindByBasicCredentials finds a sheet by basic auth credentials (auth_type = 'basic')
func (r *allowedSheetRepo) FindByBasicCredentials(ctx context.Context, username, password string) (models.AllowedSheet, error) {
	var sheet models.AllowedSheet
	err := r.db.GetContext(ctx, &sheet, `
		SELECT * FROM allowed_sheets 
		WHERE auth_type = 'basic' AND auth_basic_username = $1 AND is_public = true
	`, username)
	if err != nil {
		return sheet, err
	}

	// Verify password hash
	if sheet.AuthBasicPasswordHash == nil {
		return models.AllowedSheet{}, ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*sheet.AuthBasicPasswordHash), []byte(password)); err != nil {
		return models.AllowedSheet{}, ErrUnauthorized
	}

	return sheet, nil
}

// UpdateAuth updates the authentication type and credentials for a sheet
func (r *allowedSheetRepo) UpdateAuth(ctx context.Context, sheetID uuid.UUID, authType string, bearerToken, basicUsername, basicPasswordHash *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE allowed_sheets 
		SET auth_type = $1,
		    auth_bearer_token = $2,
		    auth_basic_username = $3,
		    auth_basic_password_hash = $4,
		    updated_at = NOW()
		WHERE id = $5
	`, authType, bearerToken, basicUsername, basicPasswordHash, sheetID)
	return err
}

func generateAPIKey() string {
	b := make([]byte, 24)
	rand.Read(b)
	return "gsheet_" + base64.URLEncoding.EncodeToString(b)[:32]
}
