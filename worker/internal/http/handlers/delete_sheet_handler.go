package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeletePublic handles DELETE /v1/:api_key?collection=assets&where={...}
func (h *SheetHandler) DeletePublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	// Parse query params
	collection := c.Query("collection")
	where := c.Query("where")

	// Find the sheet by API key
	sheet, err := h.sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid api key or sheet not found"})
		return
	}

	// Check if DELETE method is allowed
	if !isMethodAllowed(sheet.AllowedMethods, "DELETE") {
		c.JSON(http.StatusForbidden, gin.H{"error": "DELETE method not enabled for this sheet"})
		return
	}

	// Get the user to access their Google tokens
	user, err := h.userRepo.FindByID(c.Request.Context(), sheet.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user credentials"})
		return
	}

	if user.GoogleAccessToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sheet owner needs to reconnect Google account"})
		return
	}

	// Determine range (use collection as sheet name, fallback to default)
	targetRange := collection
	if targetRange == "" && sheet.DefaultRange != nil {
		targetRange = *sheet.DefaultRange
	}
	if targetRange == "" {
		targetRange = "Sheet1"
	}

	// Fetch current sheet data
	sheetData, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange)
	if err != nil || len(sheetData) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet data", "details": err.Error()})
		return
	}
	rows := transformToJSON(sheetData)

	// Filter rows to delete
	var cond map[string]interface{}
	if where != "" {
		if err := json.Unmarshal([]byte(where), &cond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid where filter", "details": err.Error()})
			return
		}
	}

	match := false
	rowIndexToDelete := -1
	for i, row := range rows {
		if where != "" {
			for k, v := range cond {
				if row[k] == v {
					match = true
					rowIndexToDelete = i + 1 // header row additional
					break
				}
			}
		}
	}

	if !match {
		c.JSON(http.StatusNotFound, gin.H{"error": "no rows matched to delete"})
		return
	}

	// Write updated data back to sheet
	if err := h.deleteSheetDataAtRowIndex(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange, int64(rowIndexToDelete)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update data", "details": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
