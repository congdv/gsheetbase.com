package repository

import (
	"context"
	"strings"
	"time"

	"gsheetbase/shared/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo interface {
	CreateOAuth(ctx context.Context, email, provider, providerID string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.User, error)
	UpdateGoogleTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiry time.Time) error
	UpdateGoogleScopes(ctx context.Context, userID uuid.UUID, scopes []string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateOAuth(ctx context.Context, email, provider, providerID string) (models.User, error) {
	now := time.Now()
	u := models.User{
		ID:         uuid.New(),
		Email:      strings.ToLower(email),
		Provider:   provider,
		ProviderID: providerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO users (id, email, provider, provider_id, created_at, updated_at)
		VALUES (:id, :email, :provider, :provider_id, :created_at, :updated_at)
	`, &u)

	return u, err
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	return u, err
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id = $1`, id)
	return u, err
}

func (r *userRepo) UpdateGoogleTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiry time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users 
		SET google_access_token = $1, 
		    google_refresh_token = $2, 
		    google_token_expiry = $3,
		    updated_at = NOW()
		WHERE id = $4
	`, accessToken, refreshToken, expiry, userID)
	return err
}

func (r *userRepo) UpdateGoogleScopes(ctx context.Context, userID uuid.UUID, scopes []string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users 
		SET google_scopes = $1,
		    updated_at = NOW()
		WHERE id = $2
	`, pq.Array(scopes), userID)
	return err
}
