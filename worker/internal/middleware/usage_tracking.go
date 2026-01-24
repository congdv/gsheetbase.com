package middleware

import (
	"context"
	"log"
	"time"

	"gsheetbase/shared/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UsageEvent represents a single API usage event
type UsageEvent struct {
	APIKey    string
	UserID    uuid.UUID
	SheetID   uuid.UUID
	Method    string
	Timestamp time.Time
}

// UsageTracker manages async usage tracking
type UsageTracker struct {
	usageRepo   repository.UsageRepo
	eventChan   chan UsageEvent
	stopChan    chan struct{}
	workerCount int
}

// NewUsageTracker creates a new usage tracker with background workers
func NewUsageTracker(usageRepo repository.UsageRepo, workerCount int) *UsageTracker {
	tracker := &UsageTracker{
		usageRepo:   usageRepo,
		eventChan:   make(chan UsageEvent, 10000),
		stopChan:    make(chan struct{}),
		workerCount: workerCount,
	}

	// Start background workers
	for i := 0; i < workerCount; i++ {
		go tracker.worker()
	}

	return tracker
}

// worker processes usage events from the channel
func (t *UsageTracker) worker() {
	for {
		select {
		case event := <-t.eventChan:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			err := t.usageRepo.IncrementDailyUsage(
				ctx,
				event.APIKey,
				event.UserID,
				event.SheetID,
				event.Timestamp,
				event.Method,
			)
			cancel()

			if err != nil {
				log.Printf("Failed to increment usage: %v", err)
			}
		case <-t.stopChan:
			return
		}
	}
}

// Track queues a usage event for async processing
func (t *UsageTracker) Track(apiKey string, userID, sheetID uuid.UUID, method string) {
	select {
	case t.eventChan <- UsageEvent{
		APIKey:    apiKey,
		UserID:    userID,
		SheetID:   sheetID,
		Method:    method,
		Timestamp: time.Now(),
	}:
	default:
		log.Printf("Usage tracking channel full, dropping event for API key: %s", apiKey)
	}
}

// Shutdown gracefully stops the usage tracker
func (t *UsageTracker) Shutdown() {
	close(t.stopChan)
	close(t.eventChan)
}

// UsageTrackingMiddleware creates a middleware that tracks API usage
func UsageTrackingMiddleware(tracker *UsageTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only track successful requests
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			apiKey := c.Param("api_key")
			method := c.Request.Method

			sheetIDRaw, sheetExists := c.Get("sheet_id")
			userIDRaw, userExists := c.Get("user_id")

			if sheetExists && userExists {
				sheetID, sheetOk := sheetIDRaw.(uuid.UUID)
				userID, userOk := userIDRaw.(uuid.UUID)

				if sheetOk && userOk {
					tracker.Track(apiKey, userID, sheetID, method)
				}
			}
		}
	}
}
