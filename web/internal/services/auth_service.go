package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"gsheetbase/shared/models"
	"gsheetbase/shared/repository"
	"gsheetbase/web/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	FindOrCreateOauthUser(ctx context.Context, email, provider, providerId string) (models.User, error)
	GenerateAccessToken(user models.User) (string, time.Time, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateRefreshToken(ctx context.Context, tokenPlaintext string) (models.User, error)
	Me(ctx context.Context, userID uuid.UUID) (models.User, error)
	UpdateGoogleTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiry time.Time) error
	UpdateGoogleScopes(ctx context.Context, userID uuid.UUID, scopes []string) error
	HasScope(ctx context.Context, userID uuid.UUID, scope string) (bool, error)
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

func (a *authService) UpdateGoogleScopes(ctx context.Context, userID uuid.UUID, scopes []string) error {
	return a.users.UpdateGoogleScopes(ctx, userID, scopes)
}

func (a *authService) HasScope(ctx context.Context, userID uuid.UUID, scope string) (bool, error) {
	user, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, s := range user.GoogleScopes {
		if s == scope {
			return true, nil
		}
	}
	return false, nil
}

// GenerateRefreshToken creates a new refresh token, hashes it, stores in DB, and returns the plaintext token
func (a *authService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate a cryptographically secure random token (32 bytes = 256 bits)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	plaintext := hex.EncodeToString(tokenBytes)

	// Hash the token before storing
	hash := sha256.Sum256([]byte(plaintext))
	hashHex := hex.EncodeToString(hash[:])

	// Calculate expiry based on config
	expiry := time.Now().AddDate(0, 0, a.cfg.JWTRefreshTTLDays)

	// Store the hash in the database
	if err := a.users.SaveRefreshToken(ctx, userID, hashHex, expiry); err != nil {
		return "", err
	}

	// Return the plaintext token (never stored, only shown once)
	return plaintext, nil
}

// ValidateRefreshToken validates a plaintext refresh token against the stored hash
func (a *authService) ValidateRefreshToken(ctx context.Context, tokenPlaintext string) (models.User, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(tokenPlaintext))
	hashHex := hex.EncodeToString(hash[:])

	// Look up user by token hash
	user, err := a.users.FindByRefreshTokenHash(ctx, hashHex)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("invalid refresh token")
		}
		return models.User{}, err
	}

	// Check if token has expired
	if user.RefreshTokenExpiry == nil || user.RefreshTokenExpiry.Before(time.Now()) {
		return models.User{}, errors.New("refresh token expired")
	}

	return user, nil
}
