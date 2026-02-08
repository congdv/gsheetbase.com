package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"gsheetbase/shared/repository"
	"gsheetbase/web/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	usageRepo repository.UsageRepo
	sheetRepo repository.AllowedSheetRepo
}

func NewAnalyticsHandler(usageRepo repository.UsageRepo, sheetRepo repository.AllowedSheetRepo) *AnalyticsHandler {
	return &AnalyticsHandler{
		usageRepo: usageRepo,
		sheetRepo: sheetRepo,
	}
}

// DailyUsageSummary represents usage stats for a single day
type DailyUsageSummary struct {
	Date       string `json:"date"`
	TotalCount int    `json:"total_count"`
	GetCount   int    `json:"get_count"`
	PostCount  int    `json:"post_count"`
	PutCount   int    `json:"put_count"`
	PatchCount int    `json:"patch_count"`
}

// GetSheetAnalytics returns usage analytics for a specific sheet
// GET /api/v1/sheets/:id/analytics?days=30
func (h *AnalyticsHandler) GetSheetAnalytics(c *gin.Context) {
	// Extract sheet ID from URL
	sheetIDStr := c.Param("id")
	sheetID, err := uuid.Parse(sheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sheet ID"})
		return
	}

	// Get authenticated user from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verify user owns this sheet
	sheet, err := h.sheetRepo.FindByID(c.Request.Context(), sheetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sheet not found"})
		return
	}

	if sheet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this sheet"})
		return
	}

	// Parse days parameter (default: 30)
	daysParam := c.DefaultQuery("days", "30")
	var days int
	if _, err := fmt.Sscanf(daysParam, "%d", &days); err != nil || days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90 // Max 90 days
	}

	// Calculate date range
	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, 0, -days)

	// Fetch usage data
	usageRecords, err := h.usageRepo.GetDailyUsageBySheet(c.Request.Context(), sheetID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch analytics"})
		return
	}

	// Aggregate by date
	dailyMap := make(map[string]*DailyUsageSummary)
	for _, record := range usageRecords {
		// Group by calendar date (midnight UTC) so counts for the same day
		// always map to the same key. Use RFC3339 timestamp at 00:00:00Z.
		d := record.RequestDate.UTC()
		dateKeyTime := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		dateKey := dateKeyTime.Format(time.RFC3339)

		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &DailyUsageSummary{
				Date: dateKey,
			}
		}

		summary := dailyMap[dateKey]
		summary.TotalCount += record.RequestCount

		switch record.Method {
		case "GET":
			summary.GetCount += record.RequestCount
		case "POST":
			summary.PostCount += record.RequestCount
		case "PUT":
			summary.PutCount += record.RequestCount
		case "PATCH":
			summary.PatchCount += record.RequestCount
		}
	}

	// Ensure every day in the range has an entry (default zeros)
	start := startDate.UTC()
	end := endDate.UTC()
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		if _, exists := dailyMap[key]; !exists {
			dailyMap[key] = &DailyUsageSummary{Date: key}
		}
	}

	// Convert map to slice
	dailyStats := make([]DailyUsageSummary, 0, len(dailyMap))
	for _, summary := range dailyMap {
		dailyStats = append(dailyStats, *summary)
	}

	sort.Slice(dailyStats, func(i, j int) bool {
		return dailyStats[i].Date < dailyStats[j].Date
	})

	c.JSON(http.StatusOK, gin.H{
		"sheet_id":    sheetID,
		"sheet_name":  sheet.SheetName,
		"api_key":     sheet.APIKey,
		"period_days": days,
		"start_date":  startDate.UTC().Format(time.RFC3339),
		"end_date":    endDate.UTC().Format(time.RFC3339),
		"daily_usage": dailyStats,
	})
}

// GetUserAnalytics returns aggregated usage analytics for all user's sheets
// GET /api/v1/analytics?days=30
func (h *AnalyticsHandler) GetUserAnalytics(c *gin.Context) {
	// Get authenticated user from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse days parameter (default: 30)
	daysParam := c.DefaultQuery("days", "30")
	var days int
	if _, err := fmt.Sscanf(daysParam, "%d", &days); err != nil || days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	// Calculate date range
	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, 0, -days)

	// Fetch usage data for all user's sheets
	usageRecords, err := h.usageRepo.GetDailyUsageByUser(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch analytics"})
		return
	}

	// Aggregate by date
	dailyMap := make(map[string]*DailyUsageSummary)
	totalRequests := 0

	for _, record := range usageRecords {
		// Group by calendar date (midnight UTC)
		d := record.RequestDate.UTC()
		dateKeyTime := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		dateKey := dateKeyTime.Format(time.RFC3339)

		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &DailyUsageSummary{
				Date: dateKey,
			}
		}

		summary := dailyMap[dateKey]
		summary.TotalCount += record.RequestCount
		totalRequests += record.RequestCount

		switch record.Method {
		case "GET":
			summary.GetCount += record.RequestCount
		case "POST":
			summary.PostCount += record.RequestCount
		case "PUT":
			summary.PutCount += record.RequestCount
		case "PATCH":
			summary.PatchCount += record.RequestCount
		}
	}

	// Ensure every day in the range has an entry (default zeros)
	start := startDate.UTC()
	end := endDate.UTC()
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		if _, exists := dailyMap[key]; !exists {
			dailyMap[key] = &DailyUsageSummary{Date: key}
		}
	}

	// Convert map to slice
	dailyStats := make([]DailyUsageSummary, 0, len(dailyMap))
	for _, summary := range dailyMap {
		dailyStats = append(dailyStats, *summary)
	}

	sort.Slice(dailyStats, func(i, j int) bool {
		return dailyStats[i].Date < dailyStats[j].Date
	})

	c.JSON(http.StatusOK, gin.H{
		"period_days":    days,
		"start_date":     startDate.UTC().Format(time.RFC3339),
		"end_date":       endDate.UTC().Format(time.RFC3339),
		"total_requests": totalRequests,
		"daily_usage":    dailyStats,
	})
}
