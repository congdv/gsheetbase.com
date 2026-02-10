package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"gsheetbase/shared/models"
	"gsheetbase/shared/repository"

	"github.com/gin-gonic/gin"
)

// SheetAuthMiddleware resolves a sheet using one of three auth methods:
// 1. API key (backward compatibility): GET /v1/:api_key
// 2. Bearer token: Authorization: Bearer <token>
// 3. Basic auth: Authorization: Basic <base64(username:password)>
//
// For auth_type = 'none', only allows access if is_public = true.
// Sets sheet_id and user_id in context for downstream handlers.
func SheetAuthMiddleware(sheetRepo repository.AllowedSheetRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Param("api_key")
		authHeader := c.GetHeader("Authorization")

		var sheet models.AllowedSheet
		var err error

		// Try API key first (backward compatibility)
		if apiKey != "" {
			sheet, err = sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
			if err == nil && sheet.AuthType == "none" {
				// API key found and is_public=true (enforced by FindByAPIKey)
				c.Set("sheet_id", sheet.ID)
				c.Set("user_id", sheet.UserID)
				c.Next()
				return
			}

			if err != nil {
				// API key provided but not found
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid api_key"})
				c.Abort()
				return
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing api_key"})
			c.Abort()
			return
		}

		// No API key: try Authorization header
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		// Parse Authorization header: "Bearer <token>" or "Basic <base64>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		scheme := parts[0]
		credentials := parts[1]

		switch scheme {
		case "Bearer":
			// Bearer token authentication
			sheet, err = sheetRepo.FindByBearerToken(c.Request.Context(), credentials)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
				c.Abort()
				return
			}

		case "Basic":
			// Basic authentication: decode base64(username:password)
			decoded, err := base64.StdEncoding.DecodeString(credentials)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Basic auth encoding"})
				c.Abort()
				return
			}

			userPass := strings.SplitN(string(decoded), ":", 2)
			if len(userPass) != 2 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Basic auth format"})
				c.Abort()
				return
			}

			username := userPass[0]
			password := userPass[1]

			sheet, err = sheetRepo.FindByBasicCredentials(c.Request.Context(), username, password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				c.Abort()
				return
			}

		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unsupported Authorization scheme; use Bearer or Basic"})
			c.Abort()
			return
		}

		// Successfully resolved sheet
		c.Set("sheet_id", sheet.ID)
		c.Set("user_id", sheet.UserID)
		c.Next()
	}
}
