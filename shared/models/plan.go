package models

// SubscriptionPlan represents the subscription tier
type SubscriptionPlan string

const (
	PlanFree       SubscriptionPlan = "free"
	PlanStarter    SubscriptionPlan = "starter"
	PlanPro        SubscriptionPlan = "pro"
	PlanEnterprise SubscriptionPlan = "enterprise"
)

// BillingPeriod represents the billing cycle
type BillingPeriod string

const (
	BillingMonthly BillingPeriod = "monthly"
	BillingAnnual  BillingPeriod = "annual"
)

// PlanLimits defines rate limits and quotas for a subscription plan
type PlanLimits struct {
	// Rate limits (per minute)
	GetRateLimit    int // GET requests per minute
	UpdateRateLimit int // POST/PUT/PATCH requests per minute

	// Daily quotas
	DailyUpdateQuota int // Maximum updates per day

	// Monthly quotas
	MonthlyGetQuota    int // Maximum GET requests per month
	MonthlyUpdateQuota int // Maximum updates per month

	// Pricing
	MonthlyPrice int // Price in cents (USD)
	AnnualPrice  int // Annual price in cents (USD)

	// Features
	CacheMinTTL     int  // Minimum cache TTL in seconds
	CustomDomain    bool // Allow custom domain
	PrioritySupport bool // Priority support access
}

// GetPlanLimits returns the limits for a given subscription plan
func GetPlanLimits(plan SubscriptionPlan) PlanLimits {
	switch plan {
	case PlanFree:
		return PlanLimits{
			GetRateLimit:       60,
			UpdateRateLimit:    2,
			DailyUpdateQuota:   200,
			MonthlyGetQuota:    10000,
			MonthlyUpdateQuota: 2000,
			MonthlyPrice:       0,
			AnnualPrice:        0,
			CacheMinTTL:        60, // Force 60s cache on free tier
			CustomDomain:       false,
			PrioritySupport:    false,
		}
	case PlanStarter:
		return PlanLimits{
			GetRateLimit:       50,
			UpdateRateLimit:    10,
			DailyUpdateQuota:   500,
			MonthlyGetQuota:    50000,
			MonthlyUpdateQuota: 5000,
			MonthlyPrice:       499,  // $4.99
			AnnualPrice:        4799, // $47.99 (~$3.99/mo)
			CacheMinTTL:        30,
			CustomDomain:       false,
			PrioritySupport:    false,
		}
	case PlanPro:
		return PlanLimits{
			GetRateLimit:       200,
			UpdateRateLimit:    50,
			DailyUpdateQuota:   2000,
			MonthlyGetQuota:    500000,
			MonthlyUpdateQuota: 50000,
			MonthlyPrice:       1999,  // $19.99
			AnnualPrice:        19199, // $191.99 (~$15.99/mo)
			CacheMinTTL:        10,
			CustomDomain:       true,
			PrioritySupport:    true,
		}
	case PlanEnterprise:
		return PlanLimits{
			GetRateLimit:       1000,   // Default, can be customized
			UpdateRateLimit:    200,    // Default, can be customized
			DailyUpdateQuota:   100000, // Default, can be customized
			MonthlyGetQuota:    10000000,
			MonthlyUpdateQuota: 1000000,
			MonthlyPrice:       9900, // $99 starting
			AnnualPrice:        0,    // Custom negotiation
			CacheMinTTL:        0,    // No minimum
			CustomDomain:       true,
			PrioritySupport:    true,
		}
	default:
		return GetPlanLimits(PlanFree)
	}
}

// GetEffectiveRateLimit returns the appropriate rate limit based on HTTP method
func (p PlanLimits) GetEffectiveRateLimit(method string) int {
	switch method {
	case "GET", "HEAD", "OPTIONS":
		return p.GetRateLimit
	case "POST", "PUT", "PATCH", "DELETE":
		return p.UpdateRateLimit
	default:
		return p.GetRateLimit
	}
}

// IsWriteMethod returns true if the method is a write operation
func IsWriteMethod(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}
