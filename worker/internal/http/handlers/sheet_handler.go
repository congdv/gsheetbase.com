package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gsheetbase/shared/repository"
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

// GetPublic handles GET /v1/:api_key
func (h *SheetHandler) GetPublic(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key is required"})
		return
	}

	// Get range from query param or use default
	rangeParam := c.Query("range")

	// Find the sheet by API key
	sheet, err := h.sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid api key or sheet not found"})
		return
	}

	// Get the user to access their Google tokens
	user, err := h.userRepo.FindByID(c.Request.Context(), sheet.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user credentials"})
		return
	}

	// Check if user has valid Google tokens
	if user.GoogleAccessToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sheet owner needs to reconnect Google account"})
		return
	}

	// Check if token is expired and we have a refresh token
	if user.GoogleTokenExpiry != nil && user.GoogleTokenExpiry.Before(time.Now()) && user.GoogleRefreshToken != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sheet owner's Google token has expired, please ask them to re-authenticate via the dashboard"})
		return
	}

	// Determine which range to fetch
	fetchRange := rangeParam
	if fetchRange == "" && sheet.DefaultRange != nil {
		fetchRange = *sheet.DefaultRange
	}
	if fetchRange == "" {
		fetchRange = "Sheet1" // Final fallback
	}

	// Fetch sheet data from Google
	data, err := h.fetchSheetData(c.Request.Context(), *user.GoogleAccessToken, sheet.SheetID, fetchRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sheet data", "details": err.Error()})
		return
	}

	// Transform data if needed
	if sheet.UseFirstRowAsHeader {
		transformed := transformToJSON(data)
		c.JSON(http.StatusOK, gin.H{"data": transformed})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
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
