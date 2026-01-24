package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PutPublic handles PUT /v1/:api_key/rows - Update/replace rows
func (h *SheetHandler) PutPublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	var req struct {
		Data  [][]interface{} `json:"data" binding:"required"`
		Range string          `json:"range" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Find the sheet by API key
	sheet, err := h.sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid api key or sheet not found"})
		return
	}

	// Check if PUT method is allowed
	if !isMethodAllowed(sheet.AllowedMethods, "PUT") {
		c.JSON(http.StatusForbidden, gin.H{"error": "PUT method not enabled for this sheet"})
		return
	}

	// Get the user to access their Google tokens
	user, err := h.userRepo.FindByID(c.Request.Context(), sheet.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user credentials"})
		return
	}

	// Set context values for usage tracking middleware
	c.Set("sheet_id", sheet.ID)
	c.Set("user_id", user.ID)

	if user.GoogleAccessToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sheet owner needs to reconnect Google account"})
		return
	}

	// Update data in sheet
	if err := h.updateSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, req.Range, req.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rows updated successfully"})
}
