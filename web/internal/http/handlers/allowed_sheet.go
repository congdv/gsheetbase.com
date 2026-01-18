package handlers

import (
	"net/http"

	"gsheetbase/shared/repository"
	"gsheetbase/web/internal/http/middleware"
	"github.com/gin-gonic/gin"
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
	DefaultRange           string `json:"default_range"`
	UseFirstRowAsHeader    bool   `json:"use_first_row_as_header"`
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
