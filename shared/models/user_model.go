package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	Email              string         `db:"email" json:"email"`
	Provider           string         `db:"provider" json:"-"`
	ProviderID         string         `db:"provider_id" json:"-"`
	GoogleAccessToken  *string        `db:"google_access_token" json:"google_access_token,omitempty"`
	GoogleRefreshToken *string        `db:"google_refresh_token" json:"-"`
	GoogleTokenExpiry  *time.Time     `db:"google_token_expiry" json:"google_token_expiry,omitempty"`
	GoogleScopes       pq.StringArray `db:"google_scopes" json:"google_scopes,omitempty"`

	// Legacy rate limit fields (deprecated, use subscription plan instead)
	RateLimitPerMinute int `db:"rate_limit_per_minute" json:"rate_limit_per_minute"`
	RateLimitBurst     int `db:"rate_limit_burst" json:"rate_limit_burst"`

	// Subscription fields
	SubscriptionPlan      SubscriptionPlan `db:"subscription_plan" json:"subscription_plan"`
	BillingPeriod         *BillingPeriod   `db:"billing_period" json:"billing_period,omitempty"`
	SubscriptionStatus    string           `db:"subscription_status" json:"subscription_status"`
	SubscriptionStartedAt *time.Time       `db:"subscription_started_at" json:"subscription_started_at,omitempty"`
	SubscriptionEndsAt    *time.Time       `db:"subscription_ends_at" json:"subscription_ends_at,omitempty"`
	StripeCustomerID      *string          `db:"stripe_customer_id" json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID  *string          `db:"stripe_subscription_id" json:"stripe_subscription_id,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// GetPlanLimits returns the subscription limits for this user
func (u *User) GetPlanLimits() PlanLimits {
	return GetPlanLimits(u.SubscriptionPlan)
}

// IsSubscriptionActive checks if the user's subscription is active
func (u *User) IsSubscriptionActive() bool {
	if u.SubscriptionPlan == PlanFree {
		return true
	}

	if u.SubscriptionStatus != "active" && u.SubscriptionStatus != "trialing" {
		return false
	}

	if u.SubscriptionEndsAt != nil && u.SubscriptionEndsAt.Before(time.Now()) {
		return false
	}

	return true
}
