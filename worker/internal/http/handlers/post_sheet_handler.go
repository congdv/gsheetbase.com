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

	headerRange := targetRange + "!1:1"

	// fetch header row
	headerData, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, headerRange)
	if err != nil || len(headerData) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet headers", "details": err.Error()})
		return
	}
	headers := headerData[0]

	// Validate the json input
	row, err := validateAndMap(headers, req.Data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate data", "details": err.Error()})
		return
	}

	fmt.Printf("Validation Error: %v\n%v\n", row, err)

	// Append the row and get the appended values from the API response
	appendResp, err := h.appendSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange, [][]interface{}{row})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to append data", "details": err.Error()})
		return
	}
	// Use the values returned in the response
	var createdRow map[string]interface{}
	if appendResp != nil && appendResp.Updates != nil && appendResp.Updates.UpdatedData != nil && len(appendResp.Updates.UpdatedData.Values) > 0 {
		createdRow = map[string]interface{}{}
		for i, h := range headers {
			headerStr, _ := h.(string)
			if i < len(appendResp.Updates.UpdatedData.Values[0]) {
				createdRow[headerStr] = appendResp.Updates.UpdatedData.Values[0][i]
			} else {
				createdRow[headerStr] = nil
			}
		}
	} else {
		// Fallback: map input row
		createdRow = map[string]interface{}{}
		for i, h := range headers {
			headerStr, _ := h.(string)
			if i < len(row) {
				createdRow[headerStr] = row[i]
			} else {
				createdRow[headerStr] = nil
			}
		}
	}

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
