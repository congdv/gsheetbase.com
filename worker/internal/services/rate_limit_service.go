package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimitService handles rate limiting logic using Redis
type RateLimitService struct {
	redis *redis.Client
}

// RateLimitResult contains the result of a rate limit check
type RateLimitResult struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   time.Time
}

// NewRateLimitService creates a new rate limit service
func NewRateLimitService(redisClient *redis.Client) *RateLimitService {
	return &RateLimitService{
		redis: redisClient,
	}
}

// CheckLimit checks if a request is allowed under the rate limit
// httpMethod is used to determine if it's a read (GET) or write (POST/PUT/PATCH) operation
// The limit parameter must be provided (from the user's subscription plan)
func (s *RateLimitService) CheckLimit(ctx context.Context, apiKey, httpMethod string, limit int) (*RateLimitResult, error) {
	// Use current minute as window, including method type for separate tracking
	now := time.Now()
	windowKey := fmt.Sprintf("rate_limit:%s:%s:%s", apiKey, httpMethod, now.Format("2006-01-02T15:04"))

	// Lua script for atomic increment and check
	script := redis.NewScript(`
		local current = redis.call('INCR', KEYS[1])
		if current == 1 then
			redis.call('EXPIRE', KEYS[1], 60)
		end
		return current
	`)

	result, err := script.Run(ctx, s.redis, []string{windowKey}).Int64()
	if err != nil {
		return nil, fmt.Errorf("rate limit check failed: %w", err)
	}

	current := int(result)
	remaining := limit - current
	if remaining < 0 {
		remaining = 0
	}

	resetAt := now.Truncate(time.Minute).Add(time.Minute)

	return &RateLimitResult{
		Allowed:   current <= limit,
		Limit:     limit,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}
