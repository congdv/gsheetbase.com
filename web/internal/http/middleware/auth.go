package middleware

import (
	"net/http"
	"strings"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/models"
	"gsheetbase/web/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ctxKey string

const (
	ctxUserID ctxKey = "userId"
	ctxUser   ctxKey = "user"
)

type accessClaims struct {
	UserId string `json:"uid"`
	jwt.RegisteredClaims
}

// Authenticate validates JWT access tokens (used for Google OAuth sessions)
func Authenticate(cfg *config.Config, authService services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Try to get token from Authorization header first, then from cookie
		var tokenStr string
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Try to get from cookie
			tokenStr, _ = ctx.Cookie("access_token")
		}

		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTAccessSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(*accessClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		uid, err := uuid.Parse(claims.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid uid"})
			return
		}

		// Fetch full user from database (including Google tokens)
		user, err := authService.Me(ctx.Request.Context(), uid)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		ctx.Set(string(ctxUserID), uid)
		ctx.Set(string(ctxUser), user)
		ctx.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx *gin.Context) (uuid.UUID, bool) {
	val, exists := ctx.Get(string(ctxUserID))
	if !exists {
		return uuid.Nil, false
	}
	uid, ok := val.(uuid.UUID)
	return uid, ok
}

// GetUserFromContext extracts full user object from context
func GetUserFromContext(ctx *gin.Context) (models.User, bool) {
	val, exists := ctx.Get(string(ctxUser))
	if !exists {
		return models.User{}, false
	}
	user, ok := val.(models.User)
	return user, ok
}
