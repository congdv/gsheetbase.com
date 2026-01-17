package handlers

import (
	"net/http"

	"github.com/congdv/go-auth/api/internal/config"
	"github.com/congdv/go-auth/api/internal/http/middleware"
	"github.com/congdv/go-auth/api/internal/services"
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
