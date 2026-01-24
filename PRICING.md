# Subscription Plans & Pricing

This document describes the subscription tiers and their limits for Gsheetbase.

## 📊 Pricing Overview

| Plan           | Monthly Price | Annual Price | GET         | UPDATE     | Daily quota       | Monthly quota                              |
| -------------- | ------------- | ------------ | ----------- | ---------- | ----------------- | ------------------------------------------ |
| **Free**       | $0            | $0           | 20 req/min  | 2 req/min  | 200 updates/day   | 2,000 updates/month<br>10,000 GETs/month   |
| **Starter**    | $4.99         | $47.99       | 50 req/min  | 10 req/min | 500 updates/day   | 50,000 GETs/month<br>5,000 updates/month   |
| **Pro**        | $19.99        | $191.99      | 200 req/min | 50 req/min | 2,000 updates/day | 500,000 GETs/month<br>50,000 updates/month |
| **Enterprise** | $99+          | Custom       | Custom      | Custom     | Custom            | Custom                                     |

### Annual Pricing Discount
- **Starter**: ~$3.99/mo when billed annually (~20% savings)
- **Pro**: ~$15.99/mo when billed annually (~20% savings)
- **Enterprise**: Negotiable

---

## 🎯 Plan Details

### Free Plan
**Perfect for:** Demos, MVPs, testing

**Limits:**
- ✅ 20 GET requests/minute
- ✅ 2 UPDATE requests/minute
- ✅ 10,000 GET requests/month
- ✅ 2,000 updates/month
- ✅ 200 updates/day

**Features:**
- Forced 60-second cache TTL (protects against abuse)
- Community support
- Standard API access

**Notes:** Very conservative limits to ensure safe first release. Prevents flooding while still allowing meaningful usage.

---

### Starter Plan
**Price:** $4.99/mo or $47.99/year  
**Perfect for:** Indie hackers, small projects, personal apps

**Limits:**
- ✅ 50 GET requests/minute
- ✅ 10 UPDATE requests/minute
- ✅ 50,000 GET requests/month
- ✅ 5,000 updates/month
- ✅ 500 updates/day

**Features:**
- Minimum 30-second cache TTL
- Email support
- Standard API access

**Notes:** Low entry price to attract early adopters. Still generous compared to competitors while remaining profitable.

---

### Pro Plan
**Price:** $19.99/mo or $191.99/year  
**Perfect for:** Production apps, agencies, growing teams

**Limits:**
- ✅ 200 GET requests/minute
- ✅ 50 UPDATE requests/minute
- ✅ 500,000 GET requests/month
- ✅ 50,000 updates/month
- ✅ 2,000 updates/day

**Features:**
- Minimum 10-second cache TTL
- Custom domain support
- Priority email support
- Advanced analytics

**Notes:** Designed for real production usage with higher quotas and faster caching.

---

### Enterprise Plan
**Price:** Starting at $99/mo (custom pricing)  
**Perfect for:** Large organizations, high-volume applications

**Limits:**
- ✅ Custom rate limits (default 1,000 GET/min, 200 UPDATE/min)
- ✅ Custom quotas (negotiable)
- ✅ Unlimited API keys
- ✅ Team access controls

**Features:**
- No minimum cache TTL
- Custom domain support
- Dedicated account manager
- SLA guarantees
- Phone & Slack support
- White-label options (if needed)

**Notes:** Limits are negotiated case-by-case. Do not publish fixed limits publicly.

---

## 🛡️ Quota Enforcement

### How It Works

1. **Per-Minute Rate Limits**
   - Separate limits for GET vs UPDATE operations
   - Enforced via Redis sliding window
   - Returns `429 Too Many Requests` when exceeded

2. **Daily Quotas** (UPDATE only)
   - Prevents excessive write operations
   - Resets at midnight UTC
   - Tracked in PostgreSQL

3. **Monthly Quotas**
   - Separate tracking for GET vs UPDATE
   - Resets on the 1st of each month
   - Enforced before processing requests

### Response Headers

Every API response includes quota information:

```
X-RateLimit-Limit: 50
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1737417660

X-Daily-Quota-Limit: 500
X-Daily-Quota-Used: 123

X-Monthly-Quota-Limit: 5000
X-Monthly-Quota-Used: 1234
```

### Error Responses

#### Rate Limit Exceeded (429)
```json
{
  "error": "Rate limit exceeded",
  "message": "You have exceeded the rate limit of 50 requests per minute for GET operations",
  "retry_after": 1737417660
}
```

#### Daily Quota Exceeded (429)
```json
{
  "error": "Daily quota exceeded",
  "message": "You have exceeded your daily quota of 500 updates. Quota resets at midnight UTC."
}
```

#### Monthly Quota Exceeded (429)
```json
{
  "error": "Monthly quota exceeded",
  "message": "You have exceeded your monthly quota of 5000 UPDATE operations. Please upgrade your plan or wait until next month.",
  "plan": "starter"
}
```

#### Subscription Inactive (402)
```json
{
  "error": "Subscription inactive",
  "message": "Your subscription is not active. Please upgrade or renew your plan."
}
```

---

## 💳 Implementation

### Database Schema

Subscription fields are stored in the `users` table:

```sql
CREATE TYPE subscription_plan AS ENUM ('free', 'starter', 'pro', 'enterprise');
CREATE TYPE billing_period AS ENUM ('monthly', 'annual');

ALTER TABLE users ADD COLUMN subscription_plan subscription_plan DEFAULT 'free';
ALTER TABLE users ADD COLUMN billing_period billing_period DEFAULT 'monthly';
ALTER TABLE users ADD COLUMN subscription_status TEXT DEFAULT 'active';
ALTER TABLE users ADD COLUMN stripe_customer_id TEXT;
ALTER TABLE users ADD COLUMN stripe_subscription_id TEXT;
```

### Plan Configuration

Plan limits are defined in [shared/models/plan.go](../shared/models/plan.go).

To change limits, update the `GetPlanLimits()` function.

### Middleware Stack

The worker API applies these middleware in order:

1. **QuotaEnforcementMiddleware** - Checks rate limits and quotas
2. **UsageTrackingMiddleware** - Tracks successful requests asynchronously

---

## 🚀 Upgrade Paths

### Free → Starter
- 2.5× more GET requests/minute
- 5× more UPDATE requests/minute
- 5× larger monthly quotas

### Starter → Pro
- 4× more GET requests/minute
- 5× more UPDATE requests/minute
- 10× larger monthly quotas
- Custom domain support
- Priority support

### Pro → Enterprise
- Custom negotiation
- Unlimited scaling potential
- Dedicated support

---

## 📈 Future Considerations

1. **Overage Charges** (optional)
   - Allow users to exceed quotas for a per-request fee
   - E.g., $0.01 per 1,000 extra requests

2. **Pay-As-You-Go Tier**
   - No monthly fee
   - $0.10 per 10,000 requests
   - Good for sporadic usage

3. **Team Plans**
   - Multi-user access
   - Shared quotas
   - Role-based permissions

4. **Add-ons**
   - Extra API keys: $5/mo each
   - Dedicated IP: $20/mo
   - Advanced analytics: $10/mo

---

## 🔧 Testing

To test quota limits locally:

```bash
# Run the test script
./scripts/test_rate_limit.sh

# Manually test with curl
curl -X GET "http://localhost:8081/v1/{API_KEY}" \
  -H "Accept: application/json"
```

Check response headers to verify quota information is returned correctly.

---

## 📝 Notes

- **Conservative Launch**: Free tier limits are intentionally low to prevent abuse during initial launch
- **Pricing Strategy**: Starter tier at $4.99 creates a low barrier to entry for paying customers
- **Enterprise Flexibility**: Do not publish Enterprise limits publicly—customize per customer
- **Cache Requirements**: Lower tiers have stricter cache requirements to reduce server load
