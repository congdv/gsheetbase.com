package handlers

import (
	"net/http"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/http/middleware"
	"gsheetbase/web/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth services.AuthService
	cfg  *config.Config
}

func NewAuthHandler(auth services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg}
}

func (h *AuthHandler) Me(c *gin.Context) {
	userId, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.auth.Me(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the access token cookie
	c.SetCookie("access_token", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
