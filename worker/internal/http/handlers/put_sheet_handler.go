package handlers

import (
	"fmt"
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
		Collection string                 `json:"collection"`
		Where      map[string]interface{} `json:"where"`
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

	// Fetch current sheet data to get headers and rows
	sheetData, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange)
	if err != nil || len(sheetData) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet data", "details": err.Error()})
		return
	}
	headers := sheetData[0]
	rows := transformToJSON(sheetData)

	// Find rows matching 'where' (if provided)
	var updatedRows []map[string]interface{}
	match := false
	rowIndex := -1
	for i, row := range rows {
		if req.Where != nil {
			for k, v := range req.Where {
				if row[k] == v {
					match = true
					rowIndex = i + 1 // header row additional
					break
				}
			}
		}
	}

	if match {
		// Use previous row values for fields not present in req.Data
		prevRow := sheetData[rowIndex]
		newRow := make([]interface{}, len(headers))
		for j, h := range headers {
			key := h
			if v, ok := req.Data[fmt.Sprintf("%v", key)]; ok {
				newRow[j] = v
			} else {
				// Use previous value if not updated
				if j < len(prevRow) {
					newRow[j] = prevRow[j]
				} else {
					newRow[j] = nil
				}
			}
		}
		// Update the row in sheetData
		sheetData[rowIndex] = newRow

		// Prepare response updated rows
		updatedRow := make(map[string]interface{})
		for j, h := range headers {
			updatedRow[fmt.Sprintf("%v", h)] = newRow[j]
		}
		updatedRows = append(updatedRows, updatedRow)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "no rows matched for update"})
		return
	}

	// Write updated data back to sheet (excluding header row)
	targetRange = targetRange + "!A2"
	if err := h.updateSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange, sheetData[1:]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update data", "details": err.Error()})
		return
	}

	// If returning fields specified, filter
	var responseRows []map[string]interface{}
	if len(req.Returning) > 0 {
		for _, row := range updatedRows {
			filtered := make(map[string]interface{})
			for _, k := range req.Returning {
				filtered[k] = row[k]
			}
			responseRows = append(responseRows, filtered)
		}
	} else {
		responseRows = updatedRows
	}

	c.JSON(http.StatusOK, gin.H{"data": responseRows})
}
