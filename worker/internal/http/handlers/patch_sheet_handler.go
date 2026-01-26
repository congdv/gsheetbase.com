package handlers

import (
	"github.com/gin-gonic/gin"
)

// PatchPublic handles PATCH /v1/:api_key/rows - Partial update
func (h *SheetHandler) PatchPublic(c *gin.Context) {
	h.updateSheetRows(c, "PATCH")
}
