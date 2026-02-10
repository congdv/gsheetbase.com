package handlers

import (
	"fmt"
	"net/http"

	"gsheetbase/shared/repository"
	"gsheetbase/web/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AllowedSheetHandler struct {
	repo repository.AllowedSheetRepo
}

func NewAllowedSheetHandler(repo repository.AllowedSheetRepo) *AllowedSheetHandler {
	return &AllowedSheetHandler{repo: repo}
}

type registerSheetRequest struct {
	SheetID     string `json:"sheet_id" binding:"required"`
	SheetName   string `json:"sheet_name"`
	Description string `json:"description"`
}

// Register adds a new sheet to the user's allowed list
func (h *AllowedSheetHandler) Register(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req registerSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet_id is required"})
		return
	}

	sheet, err := h.repo.Register(c.Request.Context(), userID, req.SheetID, req.SheetName, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register sheet", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"sheet": sheet})
}

// List returns all sheets registered by the user
func (h *AllowedSheetHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheets, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sheets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sheets": sheets})
}

// Delete removes a sheet from the user's allowed list
func (h *AllowedSheetHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("sheet_id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet_id is required"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, sheetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete sheet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sheet removed from allowed list"})
}

type publishSheetRequest struct {
	DefaultRange        string `json:"default_range"`
	UseFirstRowAsHeader bool   `json:"use_first_row_as_header"`
}

// Publish generates an API key and makes the sheet publicly accessible
func (h *AllowedSheetHandler) Publish(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	var req publishSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use defaults if not provided
		req.DefaultRange = "Sheet1"
		req.UseFirstRowAsHeader = true
	}

	// Verify the sheet belongs to the user
	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Generate API key and publish
	apiKey, err := h.repo.Publish(c.Request.Context(), sheet.ID, req.DefaultRange, req.UseFirstRowAsHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish sheet", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sheet published successfully",
		"api_key": apiKey,
	})
}

// Unpublish revokes the API key and makes the sheet private
func (h *AllowedSheetHandler) Unpublish(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	// Verify the sheet belongs to the user
	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Unpublish the sheet
	if err := h.repo.Unpublish(c.Request.Context(), sheet.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish sheet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sheet unpublished successfully"})
}

type updateWriteSettingsRequest struct {
	AllowWrite     *bool    `json:"allow_write"`
	AllowedMethods []string `json:"allowed_methods"`
}

// UpdateWriteSettings enables/disables write operations for a sheet
func (h *AllowedSheetHandler) UpdateWriteSettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	var req updateWriteSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Verify the sheet belongs to the user
	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Update write settings if provided
	if req.AllowWrite != nil {
		if err := h.repo.UpdateWriteSettings(c.Request.Context(), sheet.ID, *req.AllowWrite); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update write settings"})
			return
		}
	}

	// Update allowed methods if provided
	if req.AllowedMethods != nil {
		if err := h.repo.UpdateAllowedMethods(c.Request.Context(), sheet.ID, req.AllowedMethods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update allowed methods"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "write settings updated successfully"})
}

// ============================================================================
// Auth Management Endpoints
// ============================================================================

type setAuthTypeRequest struct {
	AuthType string `json:"auth_type" binding:"required,oneof=none bearer basic"`
}

type authStatusResponse struct {
	AuthType             string `json:"auth_type"`
	AuthBearerTokenSet   bool   `json:"auth_bearer_token_set"`
	AuthBasicUsernameSet bool   `json:"auth_basic_username_set"`
	AuthBasicPasswordSet bool   `json:"auth_basic_password_set"`
}

// GetAuthStatus returns the current auth configuration without exposing sensitive data
func (h *AllowedSheetHandler) GetAuthStatus(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	status := authStatusResponse{
		AuthType:             sheet.AuthType,
		AuthBearerTokenSet:   sheet.AuthBearerToken != nil,
		AuthBasicUsernameSet: sheet.AuthBasicUsername != nil,
		AuthBasicPasswordSet: sheet.AuthBasicPasswordHash != nil,
	}

	c.JSON(http.StatusOK, gin.H{"auth": status})
}

// SetAuthType sets the authentication type (none, bearer, basic)
// When changing type, clears old credentials
func (h *AllowedSheetHandler) SetAuthType(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	var req setAuthTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth_type is required (none or basic). Use GenerateBearerToken for bearer auth."})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Update auth type and clear old credentials
	err = h.repo.UpdateAuth(c.Request.Context(), sheet.ID, req.AuthType, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update auth type"})
		fmt.Printf("%v ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "auth type updated successfully",
		"auth_type": req.AuthType,
	})
}

type generateBearerTokenRequest struct {
	Description string `json:"description"` // optional description/label for the token
}

type generateBearerTokenResponse struct {
	TokenType string `json:"token_type"` // always "Bearer"
	Token     string `json:"token"`      // the actual token - only shown once
	ExpiresIn int    `json:"expires_in"` // -1 for no expiry
}

// GenerateBearerToken creates a new bearer token and stores it
func (h *AllowedSheetHandler) GenerateBearerToken(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Generate a new bearer token
	token := repository.GenerateBearerToken()

	// Update sheet with new bearer token (auth_type stays unchanged or gets set to bearer)
	authType := "bearer"
	err = h.repo.UpdateAuth(c.Request.Context(), sheet.ID, authType, &token, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate bearer token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": generateBearerTokenResponse{
			TokenType: "Bearer",
			Token:     token,
			ExpiresIn: -1, // no expiry for now
		},
		"message": "Bearer token generated successfully (save it securely, it won't be shown again)",
	})
}

type setBasicAuthRequest struct {
	Username string `json:"username" binding:"required,min=1,max=255"`
	Password string `json:"password" binding:"required,min=1,max=255"`
}

type setBasicAuthResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// SetBasicAuth sets basic authentication credentials (username and password)
func (h *AllowedSheetHandler) SetBasicAuth(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	var req setBasicAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	passwordHash := string(hashedPassword)

	// Update sheet with basic auth credentials
	authType := "basic"
	err = h.repo.UpdateAuth(c.Request.Context(), sheet.ID, authType, nil, &req.Username, &passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set basic auth credentials"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"auth": setBasicAuthResponse{
			Username: req.Username,
			Message:  "Basic auth credentials set successfully",
		},
	})
}

type rotateTokenRequest struct {
	// optional: if true, keep the old token active for 24h (for graceful migration)
	// for now, rotation is immediate
}

// RotateBearerToken generates a new bearer token and replaces the old one
func (h *AllowedSheetHandler) RotateBearerToken(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Verify current auth type is bearer
	if sheet.AuthType != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current auth type is not bearer"})
		return
	}

	// Generate a new bearer token
	newToken := repository.GenerateBearerToken()

	// Update with new token
	authType := "bearer"
	err = h.repo.UpdateAuth(c.Request.Context(), sheet.ID, authType, &newToken, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate bearer token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": generateBearerTokenResponse{
			TokenType: "Bearer",
			Token:     newToken,
			ExpiresIn: -1,
		},
		"message": "Bearer token rotated successfully (old token is now invalid)",
	})
}

type disableAuthRequest struct{}

// DisableAuth removes all authentication (sets auth_type to none)
func (h *AllowedSheetHandler) DisableAuth(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sheetID := c.Param("id")
	if sheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sheet id is required"})
		return
	}

	sheet, err := h.repo.FindByID(c.Request.Context(), middleware.MustParseUUID(sheetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Clear all auth credentials
	err = h.repo.UpdateAuth(c.Request.Context(), sheet.ID, "none", nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable auth"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "authentication disabled (sheet now requires is_public=true to access)"})
}

// ============================================================================
// Helper Functions
// ============================================================================

// (Token generation moved to shared/repository.GenerateBearerToken)
