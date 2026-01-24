# Rate Limiting & Usage Tracking

This document explains the rate limiting and usage tracking features implemented in Gsheetbase.

## Overview

- **Rate Limiting**: Prevents API abuse by limiting requests per API key per minute
- **Usage Tracking**: Records daily API usage statistics for analytics and billing

## Architecture

### Components

1. **Redis** - Stores rate limit counters (in-memory, fast)
2. **PostgreSQL** - Persists daily usage aggregates
3. **Rate Limit Middleware** - Checks limits before processing requests
4. **Usage Tracking Middleware** - Asynchronously logs successful requests

### Request Flow

```
Client Request
    ↓
Rate Limit Check (Redis)
    ↓ (if allowed)
Handler (Process Request)
    ↓ (if successful)
Usage Tracking (Async)
    ↓
PostgreSQL (Daily Aggregate)
```

## Configuration

### Environment Variables

Add to your `.env` file:

```bash
# Redis connection
REDIS_URL=redis://localhost:6379

# Rate limit settings (per API key)
RATE_LIMIT_PER_MINUTE=60     # Max requests per minute
RATE_LIMIT_BURST=100         # Not currently used (reserved)
USAGE_TRACK_WORKERS=3        # Background workers for async tracking
```

### Database Migration

Run the migration to create the usage tracking table:

```bash
# Apply migration (if using a migration tool)
psql $DATABASE_URL < migrations/20260120000002_add_usage_tracking.sql
```

The migration creates:
- `api_usage_daily` table for storing daily usage stats
- Indexes for efficient querying
- Rate limit fields on `users` and `allowed_sheets` tables

## Rate Limiting

### How It Works

- Uses **sliding window** algorithm
- Redis key format: `rate_limit:{api_key}:{YYYY-MM-DDTHH:MM}`
- Keys expire after 60 seconds
- Atomic increment prevents race conditions

### Response Headers

Every API response includes rate limit information:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1737417660
```

### When Limit Exceeded

HTTP Status: `429 Too Many Requests`

```json
{
  "error": "Rate limit exceeded",
  "message": "You have exceeded the rate limit of 60 requests per minute",
  "retry_after": 1737417660
}
```

### Per-User Rate Limits

The system supports customized rate limits:

1. **Default**: Set via `RATE_LIMIT_PER_MINUTE` environment variable
2. **Per-User Override**: Use `users.rate_limit_per_minute` column
3. **Per-Sheet Override**: Use `allowed_sheets.rate_limit_override` column

Priority: Sheet Override > User Limit > Default

## Usage Tracking

### Data Model

```sql
CREATE TABLE api_usage_daily (
    id UUID PRIMARY KEY,
    api_key TEXT NOT NULL,
    user_id UUID NOT NULL,
    sheet_id UUID NOT NULL,
    request_date DATE NOT NULL,
    method TEXT NOT NULL,        -- GET, POST, PUT, PATCH
    request_count INT DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    UNIQUE (api_key, request_date, method)
);
```

### How It Works

1. Middleware intercepts successful requests (200-299 status codes)
2. Queues event to buffered channel (non-blocking, 10k buffer)
3. Background workers process events asynchronously
4. Uses `INSERT ... ON CONFLICT` for atomic increments
5. Aggregates by API key, date, and HTTP method

### Benefits

- **No latency impact** - async processing via channel
- **High throughput** - buffered, parallel workers
- **Atomic counters** - PostgreSQL upsert prevents race conditions
- **Graceful degradation** - dropped events logged but don't block requests

## Analytics API

### Get Sheet Analytics

Retrieve usage stats for a specific sheet:

```bash
GET /api/sheets/:id/analytics?days=30
Authorization: Bearer <jwt_token>
```

**Response:**

```json
{
  "sheet_id": "uuid",
  "sheet_name": "My Sheet",
  "api_key": "gsheet_abc123",
  "period_days": 30,
  "start_date": "2026-01-01",
  "end_date": "2026-01-30",
  "daily_usage": [
    {
      "date": "2026-01-20",
      "total_count": 150,
      "get_count": 120,
      "post_count": 20,
      "put_count": 10,
      "patch_count": 0
    }
  ]
}
```

### Get User Analytics

Aggregate usage across all user's sheets:

```bash
GET /api/analytics?days=30
Authorization: Bearer <jwt_token>
```

**Response:**

```json
{
  "period_days": 30,
  "start_date": "2026-01-01",
  "end_date": "2026-01-30",
  "total_requests": 5420,
  "daily_usage": [
    {
      "date": "2026-01-20",
      "total_count": 240,
      "get_count": 200,
      "post_count": 30,
      "put_count": 10,
      "patch_count": 0
    }
  ]
}
```

## Production Considerations

### Redis Setup

**Railway:**

1. Add Redis plugin to your project
2. Copy `REDIS_URL` to environment variables
3. Restart worker service

**Docker Compose:**

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

### Scaling

**Worker Instances:**
- Rate limiting works across multiple workers (shared Redis)
- Usage tracking is eventually consistent
- Each instance has independent workers and channels

**Database:**
- Daily aggregates keep table size manageable
- Consider partitioning `api_usage_daily` by month for 1M+ records
- Retention policy: Archive data older than 90 days

### Monitoring

**Key Metrics:**

1. **Rate limit hits** - Track 429 responses
2. **Usage tracking lag** - Monitor channel buffer size
3. **Redis availability** - Fallback behavior on Redis failure
4. **Database write errors** - Check worker logs

**Redis Keys:**

```bash
# View current rate limits
redis-cli --scan --pattern "rate_limit:*"

# Check usage counters
redis-cli --scan --pattern "usage:*"
```

### Error Handling

**Rate Limit Check Fails:**
- Request proceeds (fail open)
- Error logged
- No rate limit headers

**Usage Tracking Fails:**
- Event dropped if channel full
- Database error logged
- Request unaffected

## Testing

### Manual Testing

```bash
# Test rate limiting
for i in {1..100}; do
  curl -X GET "http://localhost:8081/v1/gsheet_your_api_key" \
    -H "Content-Type: application/json"
done

# Check headers
curl -v "http://localhost:8081/v1/gsheet_your_api_key"

# Verify analytics
curl "http://localhost:8080/api/analytics?days=7" \
  -H "Authorization: Bearer $JWT_TOKEN"
```

### Load Testing

```bash
# Using wrk
wrk -t4 -c100 -d30s http://localhost:8081/v1/gsheet_your_api_key

# Using ab
ab -n 1000 -c 10 http://localhost:8081/v1/gsheet_your_api_key
```

## Future Enhancements

### Planned

1. **Caching Layer** - Cache Google Sheets API responses in Redis
2. **Rate Limit Tiers** - Free, Pro, Enterprise plans
3. **Real-time Analytics** - WebSocket streaming of usage metrics
4. **Alerting** - Email notifications on limit approaching
5. **Custom Time Windows** - Hourly, daily rate limits

### Possible Optimizations

1. **Batch DB Writes** - Flush usage to DB every N minutes instead of per-request
2. **Token Bucket** - Replace sliding window with token bucket for smoother bursts
3. **Distributed Rate Limiting** - Redis Cluster for high-scale deployments
4. **Write-Through Cache** - Update Redis usage counters, sync to DB periodically

## Troubleshooting

### Rate Limit Not Working

1. Check Redis connection: `redis-cli ping`
2. Verify environment variables loaded
3. Check worker logs for errors
4. Ensure middleware is applied to routes

### Usage Not Recording

1. Check if requests succeed (200-299 status)
2. Verify context values set in handlers (`sheet_id`, `user_id`)
3. Check worker logs for DB errors
4. Query database: `SELECT * FROM api_usage_daily ORDER BY created_at DESC LIMIT 10;`

### High Memory Usage

1. Reduce `USAGE_TRACK_WORKERS` count
2. Lower channel buffer size (edit middleware)
3. Check for Redis memory leaks
4. Monitor PostgreSQL connection pool

---

**Questions?** Check the main [README.md](../README.md) or open an issue on GitHub.
