package middleware

import (
	"fmt"
	"net/http"

	"gsheetbase/worker/internal/services"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware creates a middleware that enforces rate limits per API key
// NOTE: This is the legacy middleware. For production, use QuotaEnforcementMiddleware instead.
func RateLimitMiddleware(rateLimitService *services.RateLimitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Param("api_key")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		httpMethod := c.Request.Method

		// Use a default rate limit of 60 requests per minute for legacy usage
		defaultLimit := 60

		// Check rate limit
		result, err := rateLimitService.CheckLimit(c.Request.Context(), apiKey, httpMethod, defaultLimit)
		if err != nil {
			// Log error but don't block request on rate limit check failure
			c.Next()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", result.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", result.ResetAt.Unix()))

		if !result.Allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("You have exceeded the rate limit of %d requests per minute", result.Limit),
				"retry_after": result.ResetAt.Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
