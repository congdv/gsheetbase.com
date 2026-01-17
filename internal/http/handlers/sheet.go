package handlers

import (
	"net/http"

	"github.com/congdv/go-auth/api/internal/http/middleware"
	"github.com/congdv/go-auth/api/internal/services"
	"github.com/gin-gonic/gin"
)

type SheetHandler struct {
	service services.SheetService
}

func NewSheetHandler(service services.SheetService) *SheetHandler {
	return &SheetHandler{service: service}
}

type getSheetRequest struct {
	SheetID string `json:"sheet_id" binding:"required"`
	Range   string `json:"range" binding:"required"`
}

// Get returns Google Sheet data directly from Google API (requires JWT auth)
//
// SECURITY: Application-level access control with spreadsheets.readonly scope
// - Users must first register sheets via POST /api/sheets/register
// - Only registered sheets can be accessed
// - Cleaner OAuth consent: "View your Google Spreadsheets" (no edit/delete warnings)
//
// To register a sheet:
//   POST /api/sheets/register
//   { "sheet_id": "...", "sheet_name": "My Sheet", "description": "..." }
//
// Extract Sheet ID from URL:
//   https://docs.google.com/spreadsheets/d/SHEET_ID_HERE/edit
func (h *SheetHandler) Get(c *gin.Context) {
	var req getSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet_id and range are required"})
		return
	}

	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Check if user has Google access token
	if user.GoogleAccessToken == nil || *user.GoogleAccessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "google access token not found, please re-authenticate"})
		return
	}

	data, err := h.service.GetSheetData(c.Request.Context(), user.ID, *user.GoogleAccessToken, req.SheetID, req.Range)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
