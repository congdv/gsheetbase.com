package handlers

import (
	"net/http"
	"time"

	"gsheetbase/shared/models"
	"gsheetbase/shared/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionHandler handles subscription and billing endpoints
type SubscriptionHandler struct {
	userRepo  repository.UserRepo
	usageRepo repository.UsageRepo
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(userRepo repository.UserRepo, usageRepo repository.UsageRepo) *SubscriptionHandler {
	return &SubscriptionHandler{
		userRepo:  userRepo,
		usageRepo: usageRepo,
	}
}

// PlanInfo contains plan details and limits
type PlanInfo struct {
	Plan               string `json:"plan"`
	BillingPeriod      string `json:"billing_period,omitempty"`
	Status             string `json:"status"`
	SubscriptionEndsAt string `json:"subscription_ends_at,omitempty"`

	// Limits
	GetRateLimit       int `json:"get_rate_limit"`
	UpdateRateLimit    int `json:"update_rate_limit"`
	DailyUpdateQuota   int `json:"daily_update_quota"`
	MonthlyGetQuota    int `json:"monthly_get_quota"`
	MonthlyUpdateQuota int `json:"monthly_update_quota"`

	// Pricing
	MonthlyPrice int `json:"monthly_price_cents"`
	AnnualPrice  int `json:"annual_price_cents"`

	// Features
	CacheMinTTL     int  `json:"cache_min_ttl_seconds"`
	CustomDomain    bool `json:"custom_domain"`
	PrioritySupport bool `json:"priority_support"`
}

// UsageInfo contains current usage statistics
type UsageInfo struct {
	Today   UsagePeriod `json:"today"`
	Month   UsagePeriod `json:"month"`
	Updated string      `json:"updated_at"`
}

// UsagePeriod represents usage for a time period
type UsagePeriod struct {
	GetRequests    int `json:"get_requests"`
	UpdateRequests int `json:"update_requests"`
}

// GetCurrentPlan returns the user's current subscription plan and limits
// GET /api/v1/subscription/plan
func (h *SubscriptionHandler) GetCurrentPlan(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	limits := user.GetPlanLimits()

	planInfo := PlanInfo{
		Plan:               string(user.SubscriptionPlan),
		Status:             user.SubscriptionStatus,
		GetRateLimit:       limits.GetRateLimit,
		UpdateRateLimit:    limits.UpdateRateLimit,
		DailyUpdateQuota:   limits.DailyUpdateQuota,
		MonthlyGetQuota:    limits.MonthlyGetQuota,
		MonthlyUpdateQuota: limits.MonthlyUpdateQuota,
		MonthlyPrice:       limits.MonthlyPrice,
		AnnualPrice:        limits.AnnualPrice,
		CacheMinTTL:        limits.CacheMinTTL,
		CustomDomain:       limits.CustomDomain,
		PrioritySupport:    limits.PrioritySupport,
	}

	if user.BillingPeriod != nil {
		planInfo.BillingPeriod = string(*user.BillingPeriod)
	}

	if user.SubscriptionEndsAt != nil {
		planInfo.SubscriptionEndsAt = user.SubscriptionEndsAt.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{
		"plan": planInfo,
	})
}

// GetCurrentUsage returns the user's current API usage statistics
// GET /api/v1/subscription/usage
func (h *SubscriptionHandler) GetCurrentUsage(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get today's usage
	todayGets, err := h.usageRepo.GetTodayUsageCount(c.Request.Context(), userID, "GET")
	if err != nil {
		todayGets = 0 // Default to 0 on error
	}

	todayUpdates, err := h.usageRepo.GetTodayUsageCount(c.Request.Context(), userID, "UPDATE")
	if err != nil {
		todayUpdates = 0
	}

	// Get monthly usage
	monthlyGets, err := h.usageRepo.GetMonthlyUsageCount(c.Request.Context(), userID, "GET")
	if err != nil {
		monthlyGets = 0
	}

	monthlyUpdates, err := h.usageRepo.GetMonthlyUsageCount(c.Request.Context(), userID, "UPDATE")
	if err != nil {
		monthlyUpdates = 0
	}

	usageInfo := UsageInfo{
		Today: UsagePeriod{
			GetRequests:    todayGets,
			UpdateRequests: todayUpdates,
		},
		Month: UsagePeriod{
			GetRequests:    monthlyGets,
			UpdateRequests: monthlyUpdates,
		},
		Updated: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{
		"usage": usageInfo,
	})
}

// GetAvailablePlans returns all available subscription plans
// GET /api/v1/subscription/plans
func (h *SubscriptionHandler) GetAvailablePlans(c *gin.Context) {
	plans := []PlanInfo{
		convertPlanToInfo(models.PlanFree),
		convertPlanToInfo(models.PlanStarter),
		convertPlanToInfo(models.PlanPro),
		convertPlanToInfo(models.PlanEnterprise),
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}

// convertPlanToInfo converts a plan to PlanInfo
func convertPlanToInfo(plan models.SubscriptionPlan) PlanInfo {
	limits := models.GetPlanLimits(plan)

	return PlanInfo{
		Plan:               string(plan),
		GetRateLimit:       limits.GetRateLimit,
		UpdateRateLimit:    limits.UpdateRateLimit,
		DailyUpdateQuota:   limits.DailyUpdateQuota,
		MonthlyGetQuota:    limits.MonthlyGetQuota,
		MonthlyUpdateQuota: limits.MonthlyUpdateQuota,
		MonthlyPrice:       limits.MonthlyPrice,
		AnnualPrice:        limits.AnnualPrice,
		CacheMinTTL:        limits.CacheMinTTL,
		CustomDomain:       limits.CustomDomain,
		PrioritySupport:    limits.PrioritySupport,
	}
}

// getUserIDFromContext extracts user ID from gin context
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDRaw, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, http.ErrNoCookie
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return uuid.Nil, http.ErrNoCookie
	}

	return userID, nil
}
