package handlers

import (
	"net/http"

	"gsheetbase/web/internal/http/middleware"
	"gsheetbase/web/internal/repository"
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
