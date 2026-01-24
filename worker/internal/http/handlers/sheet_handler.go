package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gsheetbase/shared/repository"

	"encoding/json"
	"sort"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetHandler struct {
	sheetRepo repository.AllowedSheetRepo
	userRepo  repository.UserRepo
}

func NewSheetHandler(sheetRepo repository.AllowedSheetRepo, userRepo repository.UserRepo) *SheetHandler {
	return &SheetHandler{
		sheetRepo: sheetRepo,
		userRepo:  userRepo,
	}
}

// isMethodAllowed checks if a specific HTTP method is allowed for the sheet
func isMethodAllowed(allowedMethods []string, method string) bool {
	for _, m := range allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// GraphQL-style GET /v1/:api_key?collection=Sheet1&fields=asset,location&where={"owner":"Homeowner"}&orderBy=asset&limit=2&offset=0
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

// Helper: filterRows applies simple equality filter from where JSON
func filterRows(rows []map[string]interface{}, where string) []map[string]interface{} {
	if where == "" {
		return rows
	}
	// Parse where as map[string]interface{}
	var cond map[string]interface{}
	err := json.Unmarshal([]byte(where), &cond)
	if err != nil {
		return rows // ignore invalid filter
	}
	filtered := make([]map[string]interface{}, 0)
	for _, row := range rows {
		match := true
		for k, v := range cond {
			if fmt.Sprintf("%v", row[k]) != fmt.Sprintf("%v", v) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

// Helper: selectFields returns only requested fields
func selectFields(rows []map[string]interface{}, fields string) []map[string]interface{} {
	if fields == "" {
		return rows
	}
	keys := strings.Split(fields, ",")
	selected := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		obj := make(map[string]interface{})
		for _, k := range keys {
			k = strings.TrimSpace(k)
			obj[k] = row[k]
		}
		selected = append(selected, obj)
	}
	return selected
}

// Helper: orderRows sorts by a single field
func orderRows(rows []map[string]interface{}, orderBy string) []map[string]interface{} {
	if orderBy == "" {
		return rows
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return fmt.Sprintf("%v", rows[i][orderBy]) < fmt.Sprintf("%v", rows[j][orderBy])
	})
	return rows
}

// Helper: paginateRows returns paginated slice and pagination info
func paginateRows(rows []map[string]interface{}, limitStr, offsetStr string) ([]map[string]interface{}, map[string]interface{}) {
	limit, err1 := strconv.Atoi(limitStr)
	offset, err2 := strconv.Atoi(offsetStr)
	if err1 != nil || err2 != nil || limit <= 0 || offset < 0 {
		return rows, nil
	}
	total := len(rows)
	end := offset + limit
	if end > total {
		end = total
	}
	paginated := rows
	if offset < total {
		paginated = rows[offset:end]
	} else {
		paginated = []map[string]interface{}{}
	}
	pagination := map[string]interface{}{
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"nextOffset": end,
	}
	return paginated, pagination
}

func (h *SheetHandler) fetchSheetData(ctx context.Context, accessToken, sheetID, rangeStr string) ([][]interface{}, error) {
	// Create OAuth2 token
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	// Create OAuth2 config (scopes don't matter for just creating a client with existing token)
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}

	client := config.Client(ctx, token)

	// Create Sheets service
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %w", err)
	}

	// Get spreadsheet data
	resp, err := srv.Spreadsheets.Values.Get(sheetID, rangeStr).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}

// transformToJSON converts a 2D array to JSON objects using first row as headers
func transformToJSON(data [][]interface{}) []map[string]interface{} {
	if len(data) == 0 {
		return []map[string]interface{}{}
	}

	// First row is the header
	headers := data[0]
	result := make([]map[string]interface{}, 0, len(data)-1)

	// Process remaining rows
	for i := 1; i < len(data); i++ {
		row := data[i]
		obj := make(map[string]interface{})

		for j, header := range headers {
			headerStr := fmt.Sprintf("%v", header)
			if j < len(row) {
				obj[headerStr] = row[j]
			} else {
				obj[headerStr] = nil
			}
		}

		result = append(result, obj)
	}

	return result
}

// PostPublic handles POST /v1/:api_key/rows - Append new rows
func (h *SheetHandler) PostPublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	var req struct {
		Data  [][]interface{} `json:"data" binding:"required"`
		Range string          `json:"range"`
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

	// Determine range
	targetRange := req.Range
	if targetRange == "" && sheet.DefaultRange != nil {
		targetRange = *sheet.DefaultRange
	}
	if targetRange == "" {
		targetRange = "Sheet1"
	}

	// Append data to sheet
	if err := h.appendSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, targetRange, req.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to append data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rows appended successfully"})
}

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

// PatchPublic handles PATCH /v1/:api_key/rows - Partial update
func (h *SheetHandler) PatchPublic(c *gin.Context) {
	// PATCH uses same logic as PUT for Google Sheets
	h.PutPublic(c)
}

func (h *SheetHandler) appendSheetData(ctx context.Context, accessToken, sheetID, rangeStr string, data [][]interface{}) error {
	token := &oauth2.Token{AccessToken: accessToken}
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}
	client := config.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create sheets service: %w", err)
	}

	valueRange := &sheets.ValueRange{
		Values: data,
	}

	_, err = srv.Spreadsheets.Values.Append(sheetID, rangeStr, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("unable to append data to sheet: %w", err)
	}

	return nil
}

func (h *SheetHandler) updateSheetData(ctx context.Context, accessToken, sheetID, rangeStr string, data [][]interface{}) error {
	token := &oauth2.Token{AccessToken: accessToken}
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}
	client := config.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create sheets service: %w", err)
	}

	valueRange := &sheets.ValueRange{
		Values: data,
	}

	_, err = srv.Spreadsheets.Values.Update(sheetID, rangeStr, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("unable to update sheet data: %w", err)
	}

	return nil
}
