package middleware

import (
	"fmt"
	"net/http"

	"gsheetbase/shared/models"
	"gsheetbase/shared/repository"
	"gsheetbase/worker/internal/services"

	"github.com/gin-gonic/gin"
)

// QuotaEnforcementMiddleware creates a middleware that enforces both rate limits and quotas
// It checks:
// 1. Per-minute rate limits (GET vs UPDATE)
// 2. Daily quotas (for UPDATE operations)
// 3. Monthly quotas (GET and UPDATE separately)
func QuotaEnforcementMiddleware(
	rateLimitService *services.RateLimitService,
	usageRepo repository.UsageRepo,
	userRepo repository.UserRepo,
	sheetRepo repository.AllowedSheetRepo,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Param("api_key")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		httpMethod := c.Request.Method

		// First, find the sheet by API key to get the user ID
		sheet, err := sheetRepo.FindByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid API key or sheet not found"})
			c.Abort()
			return
		}

		// Set context values for downstream handlers
		c.Set("sheet_id", sheet.ID)
		c.Set("user_id", sheet.UserID)

		// Get user from database to check subscription plan
		user, err := userRepo.FindByID(c.Request.Context(), sheet.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			c.Abort()
			return
		}

		// Check if subscription is active
		if !user.IsSubscriptionActive() {
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error":   "Subscription inactive",
				"message": "Your subscription is not active. Please upgrade or renew your plan.",
			})
			c.Abort()
			return
		}

		// Get plan limits
		planLimits := user.GetPlanLimits()

		// Determine if this is a write operation
		isWrite := models.IsWriteMethod(httpMethod)
		methodCategory := "GET"
		if isWrite {
			methodCategory = "UPDATE"
		}

		// --- 1. Check per-minute rate limit ---
		effectiveRateLimit := planLimits.GetEffectiveRateLimit(httpMethod)
		rateLimitResult, err := rateLimitService.CheckLimit(c.Request.Context(), apiKey, httpMethod, effectiveRateLimit)
		if err != nil {
			// Log error but don't block request on rate limit check failure
			c.Next()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimitResult.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rateLimitResult.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", rateLimitResult.ResetAt.Unix()))

		if !rateLimitResult.Allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("You have exceeded the rate limit of %d requests per minute for %s operations", rateLimitResult.Limit, methodCategory),
				"retry_after": rateLimitResult.ResetAt.Unix(),
			})
			c.Abort()
			return
		}

		// --- 2. Check daily quota (only for write operations) ---
		if isWrite && planLimits.DailyUpdateQuota > 0 {
			dailyCount, err := usageRepo.GetTodayUsageCount(c.Request.Context(), sheet.UserID, methodCategory)
			if err != nil {
				// Log but don't block
				c.Next()
				return
			}

			c.Header("X-Daily-Quota-Limit", fmt.Sprintf("%d", planLimits.DailyUpdateQuota))
			c.Header("X-Daily-Quota-Used", fmt.Sprintf("%d", dailyCount))

			if dailyCount >= planLimits.DailyUpdateQuota {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "Daily quota exceeded",
					"message": fmt.Sprintf("You have exceeded your daily quota of %d updates. Quota resets at midnight UTC.", planLimits.DailyUpdateQuota),
				})
				c.Abort()
				return
			}
		}

		// --- 3. Check monthly quota ---
		var monthlyQuota int
		if isWrite {
			monthlyQuota = planLimits.MonthlyUpdateQuota
		} else {
			monthlyQuota = planLimits.MonthlyGetQuota
		}

		if monthlyQuota > 0 {
			monthlyCount, err := usageRepo.GetMonthlyUsageCount(c.Request.Context(), sheet.UserID, methodCategory)
			if err != nil {
				// Log but don't block
				c.Next()
				return
			}

			c.Header("X-Monthly-Quota-Limit", fmt.Sprintf("%d", monthlyQuota))
			c.Header("X-Monthly-Quota-Used", fmt.Sprintf("%d", monthlyCount))

			if monthlyCount >= monthlyQuota {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "Monthly quota exceeded",
					"message": fmt.Sprintf("You have exceeded your monthly quota of %d %s operations. Please upgrade your plan or wait until next month.", monthlyQuota, methodCategory),
					"plan":    string(user.SubscriptionPlan),
				})
				c.Abort()
				return
			}
		}

		// All checks passed
		c.Next()
	}
}
