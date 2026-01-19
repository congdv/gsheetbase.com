# Gsheetbase Worker API

Public API service for accessing Google Sheets data via API keys.

## Overview

The worker service provides a public REST API endpoint that allows anonymous access to Google Sheets data using API keys. This is separate from the main web service for better security isolation and independent scaling.

## Architecture

- **Web Service** (port 8080): Private admin API with JWT authentication
  - User management, OAuth flow
  - Sheet registration and management
  - API key generation (publish/unpublish endpoints)

- **Worker Service** (port 8081): Public data API with API key authentication
  - Anonymous sheet data access
  - Data transformation (2D arrays → JSON objects)
  - Rate limiting and caching (future)

## API Endpoint

```
GET /v1/:api_key
```

### Query Parameters

- `range` (optional): Override the default sheet range (e.g., `?range=A1:Z100`)
- `nocache` (optional): Bypass cache (future feature)

### Response Format

**With `use_first_row_as_header: true` (default):**
```json
{
  "data": [
    {"Name": "Alice", "Age": 30, "City": "NYC"},
    {"Name": "Bob", "Age": 25, "City": "SF"}
  ]
}
```

**With `use_first_row_as_header: false`:**
```json
{
  "data": [
    ["Name", "Age", "City"],
    ["Alice", 30, "NYC"],
    ["Bob", 25, "SF"]
  ]
}
```

## Development

### Run locally
```bash
make run-worker
```

### Build
```bash
make build-worker
```

### Environment Variables
```
PORT=8081
DATABASE_URL=postgres://user:pass@localhost:5432/gsheetbase
```

## Deployment (Railway)

Deploy as a separate service:

1. Create new Railway service
2. Set root directory to `/`
3. Build command: `go build -o bin/worker worker/cmd/api/main.go`
4. Start command: `./bin/worker`
5. Environment: `PORT`, `DATABASE_URL`

## CORS Configuration

The worker API has permissive CORS settings for public access:
- AllowOrigins: `*`
- AllowMethods: `GET`, `OPTIONS`
- AllowCredentials: `false`

## Future Enhancements

- [ ] Redis caching with configurable TTL
- [ ] Rate limiting per API key (Redis-based)
- [ ] Response format options (CSV, XML)
- [ ] Field filtering
- [ ] Pagination for large datasets
- [ ] API usage analytics
