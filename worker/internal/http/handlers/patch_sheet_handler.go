package handlers

import (
	"github.com/gin-gonic/gin"
)

// PatchPublic handles PATCH /v1/:api_key/rows - Partial update
func (h *SheetHandler) PatchPublic(c *gin.Context) {
	// PATCH uses same logic as PUT for Google Sheets
	h.PutPublic(c)
}
