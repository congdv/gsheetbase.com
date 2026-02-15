package handlers

import (
	"fmt"
	"net/http"
	"time"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/services"

	"github.com/gin-gonic/gin"
)

type RefreshAuthHandler struct {
	auth services.AuthService
	cfg  *config.Config
}

func NewRefreshAuthHandler(auth services.AuthService, cfg *config.Config) *RefreshAuthHandler {
	return &RefreshAuthHandler{auth: auth, cfg: cfg}
}

// RefreshSession validates the refresh token cookie and issues a new access token
func (h *RefreshAuthHandler) RefreshSession(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	// Validate the refresh token and get the user
	user, err := h.auth.ValidateRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Generate a new access token
	accessToken, exp, err := h.auth.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue new access token"})
		return
	}

	// Set the new access token as HTTP-only cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, int(time.Until(exp).Seconds()), "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	fmt.Printf("DEBUG: Session refreshed for user %s\n", user.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"message": "session refreshed successfully",
		"user":    user,
	})
}
