package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/models"
	"gsheetbase/web/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	FindOrCreateOauthUser(ctx context.Context, email, provider, providerId string) (models.User, error)
	GenerateAccessToken(user models.User) (string, time.Time, error)
	Me(ctx context.Context, userID uuid.UUID) (models.User, error)
	UpdateGoogleTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiry time.Time) error
}

type authService struct {
	cfg   *config.Config
	users repository.UserRepo
}

func NewAuthService(cfg *config.Config, users repository.UserRepo) AuthService {
	return &authService{
		cfg:   cfg,
		users: users,
	}
}

func (a *authService) FindOrCreateOauthUser(ctx context.Context, email, provider, providerId string) (models.User, error) {
	u, err := a.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u, err = a.users.CreateOAuth(ctx, email, provider, providerId)
			if err != nil {
				return models.User{}, err
			}
		} else {
			return models.User{}, err
		}
	}
	return u, nil
}

func (a *authService) Me(ctx context.Context, userID uuid.UUID) (models.User, error) {
	return a.users.FindByID(ctx, userID)
}

type accessClaims struct {
	UserId string `json:"uid"`
	jwt.RegisteredClaims
}

func (a *authService) GenerateAccessToken(user models.User) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(time.Duration(a.cfg.JWTAccessTTLMinutes) * time.Minute)

	claims := accessClaims{
		UserId: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(a.cfg.JWTAccessSecret))
	return s, exp, err
}

func (a *authService) UpdateGoogleTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiry time.Time) error {
	return a.users.UpdateGoogleTokens(ctx, userID, accessToken, refreshToken, expiry)
}
