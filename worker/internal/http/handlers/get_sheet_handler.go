package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPublic handles GET /v1/:api_key?collection=Sheet1&fields=asset,location&where={"owner":"Homeowner"}&orderBy=asset&limit=2&offset=0
func (h *SheetHandler) GetPublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	// Parse query params
	collection := c.Query("collection")
	fields := c.Query("fields")
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")
	orderBy := c.Query("orderBy")
	where := c.Query("where")

	// Find the sheet by API key
	sheet, err := h.sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid api key or sheet not found"})
		return
	}

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

	// Determine range
	fetchRange := collection
	if fetchRange == "" && sheet.DefaultRange != nil {
		fetchRange = *sheet.DefaultRange
	}
	if fetchRange == "" {
		fetchRange = "Sheet1"
	}

	// Fetch sheet data
	data, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, fetchRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet data", "details": err.Error()})
		return
	}

	// Transform to JSON
	rows := transformToJSON(data)

	// Apply simple where filter (equality only)
	filtered := filterRows(rows, where)

	// Select fields
	selected := selectFields(filtered, fields)

	// Order by
	ordered := orderRows(selected, orderBy)

	// Pagination only if explicitly set
	paginated, pagination := paginateRows(ordered, limit, offset)
	hasExplicitPagination := c.Query("limit") != "" || c.Query("offset") != ""
	if hasExplicitPagination && pagination != nil {
		c.JSON(http.StatusOK, gin.H{"data": paginated, "pagination": pagination})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": paginated})
	}
}
