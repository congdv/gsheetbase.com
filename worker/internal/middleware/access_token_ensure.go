package middleware

import (
	"net/http"
	"time"

	"gsheetbase/shared/repository"
	"gsheetbase/worker/internal/services"

	"github.com/gin-gonic/gin"
)

// AccessTokenEnsureMiddleware ensures a valid Google access token for API requests using api_key
func AccessTokenEnsureMiddleware(sheetRepo repository.AllowedSheetRepo, userRepo repository.UserRepo, authService *services.AuthService, clientId, clientSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Param("api_key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing api_key"})
			return
		}

		ctx := c.Request.Context()
		sheet, err := sheetRepo.FindByAPIKey(ctx, apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api_key"})
			return
		}

		user, err := userRepo.FindByID(ctx, sheet.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		accessToken := ""
		refreshToken := ""
		// Use global clientId/clientSecret from config
		expiry := time.Time{}

		if user.GoogleAccessToken != nil {
			accessToken = *user.GoogleAccessToken
		}
		if user.GoogleRefreshToken != nil {
			refreshToken = *user.GoogleRefreshToken
		}
		if user.GoogleTokenExpiry != nil {
			expiry = *user.GoogleTokenExpiry
		}
		// If you store clientId/clientSecret per user, add here

		// Renew if expired or expiring in <5 min
		if expiry.IsZero() || expiry.Before(time.Now().Add(5*time.Minute)) {
			resp, err := authService.RenewAccessToken(ctx, clientId, clientSecret, refreshToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token renewal failed"})
				return
			}
			accessToken = resp.AccessToken
			expiry = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
			// Update user tokens in DB (refreshToken may be empty)
			_ = userRepo.UpdateGoogleTokens(ctx, user.ID, accessToken, refreshToken, expiry)
		}

		// Attach access token to context for downstream use
		c.Set("google_access_token", accessToken)
		c.Next()
	}
}
