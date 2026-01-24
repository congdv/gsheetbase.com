package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostPublic handles POST /v1/:api_key/rows - Append new rows
func (h *SheetHandler) PostPublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	var req struct {
		Collection string                 `json:"collection" binding:"required"`
		Data       map[string]interface{} `json:"data" binding:"required"`
		Returning  []string               `json:"returning"`
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

	// Check if POST method is allowed
	if !isMethodAllowed(sheet.AllowedMethods, "POST") {
		c.JSON(http.StatusForbidden, gin.H{"error": "POST method not enabled for this sheet"})
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

	// Determine range (use collection as sheet name, fallback to default)
	targetRange := req.Collection
	if targetRange == "" && sheet.DefaultRange != nil {
		targetRange = *sheet.DefaultRange
	}
	if targetRange == "" {
		targetRange = "Sheet1"
	}

	// Fetch current sheet data to get headers
	sheetData, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange)
	if err != nil || len(sheetData) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet headers", "details": err.Error()})
		return
	}
	headers := sheetData[0]

	// Map req.Data to row in header order
	row := make([]interface{}, len(headers))
	for i, h := range headers {
		key := h
		if v, ok := req.Data[fmt.Sprintf("%v", key)]; ok {
			row[i] = v
		} else {
			row[i] = nil
		}
	}

	// Append the row
	if err := h.appendSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange, [][]interface{}{row}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to append data", "details": err.Error()})
		return
	}

	// Fetch the updated sheet data to get the last row (assume append at end)
	updatedData, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange)
	if err != nil || len(updatedData) <= 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated sheet data", "details": err.Error()})
		return
	}
	jsonRows := transformToJSON(updatedData)
	createdRow := jsonRows[len(jsonRows)-1]

	// If returning fields specified, filter
	if len(req.Returning) > 0 {
		filtered := make(map[string]interface{})
		for _, k := range req.Returning {
			filtered[k] = createdRow[k]
		}
		createdRow = filtered
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdRow})
}
