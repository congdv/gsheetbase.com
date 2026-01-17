# Google Sheets API Backend

A Go backend service that authenticates with Google OAuth and provides secure, application-level access control to Google Sheets.

## Features

- ✅ Google OAuth authentication with `spreadsheets.readonly` scope
- ✅ Application-level per-sheet access control
- ✅ JWT tokens for authenticated user sessions
- ✅ Direct Google Sheets API access (no data storage)
- ✅ Cleaner OAuth consent screen (read-only, no scary warnings)

## Security Model

### OAuth Scope: `spreadsheets.readonly`

This app uses the `spreadsheets.readonly` OAuth scope with **application-level access control**:

- ✅ Clean OAuth consent: "View your Google Spreadsheets" (no edit/delete warnings)
- ✅ Users explicitly register specific sheets they want to access
- ✅ Only registered sheets can be accessed via the API
- ✅ Database tracks which sheets each user has authorized
- ✅ Read-only access enforced by OAuth scope

### How It Works

1. User authenticates with Google OAuth (`spreadsheets.readonly` scope)
2. User registers specific sheets they want to use:
   ```bash
   POST /api/sheets/register
   { "sheet_id": "1BxiMVs...", "sheet_name": "My Data", "description": "Sales data" }
   ```
3. App stores sheet registration in database
4. User can now access registered sheets via `/api/sheets/data`
5. Attempting to access non-registered sheets returns 403 error

## Architecture

### Tech Stack
- **Framework**: Gin
- **Database**: PostgreSQL (minimal - only stores users)
- **Authentication**: Google OAuth 2.0

### Database Schema

#### Users
- OAuth-only authentication (no passwords)
- Stores Google provider info

## API Endpoints

### Authentication (Google OAuth)

#### `GET /api/auth/google/start`
Initiates Google OAuth flow

#### `GET /api/auth/google/callback`
OAuth callback handler, returns JWT token

#### `GET /api/auth/me`
Get current user info (requires JWT)

### Sheet Registration (Requires JWT)

#### `POST /api/sheets/register`
Register a new sheet for access
```json
{
  "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
  "sheet_name": "Sales Data 2026",
  "description": "Monthly sales tracking"
}
```

Returns:
```json
{
  "sheet": {
    "id": "...",
    "user_id": "...",
    "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "sheet_name": "Sales Data 2026",
    "description": "Monthly sales tracking",
    "created_at": "2026-01-17T..."
  }
}
```

#### `GET /api/sheets/registered`
List all sheets registered by the user

#### `DELETE /api/sheets/registered/:sheet_id`
Remove a sheet from registered list

### Sheet Data Access (Requires JWT + Registration)

#### `POST /api/sheets/data`
Read data from a registered sheet
```json
{
  "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
  "range": "Sheet1!A1:Z100"
}
```

Returns:
```json
{
  "data": [
    ["Header1", "Header2", "Header3"],
    ["Value1", "Value2", "Value3"]
  ]
}
```

**Note**: Sheet must be registered first via `/api/sheets/register`

## Setup

### Prerequisites
- Go 1.23+
- PostgreSQL
- Google Cloud Console project with OAuth credentials

### Environment Variables

Create a `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/gsheetbase?sslmode=disable

# JWT for session management
JWT_ACCESS_SECRET=your-secret-key
JWT_ACCESS_TTL_MINUTES=60

# CORS
FRONTEND_ORIGIN=http://localhost:3000
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

### Google Cloud Console Setup

1. Create a project in [Google Cloud Console](https://console.cloud.google.com)
2. Enable **Google Sheets API**
3. Create OAuth 2.0 credentials
4. Add authorized redirect URIs:
   - `http://localhost:8080/api/auth/google/callback`
5. Set OAuth consent screen:
   - Scopes: `spreadsheets.readonly` ("View your Google Spreadsheets")

### Database Setup

Install dbmate:
```bash
# macOS
brew install dbmate

# Linux
curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
chmod +x /usr/local/bin/dbmate
```

Run migrations:
```bash
make migrate-up
# or
dbmate up
```

Rollback:
```bash
make migrate-down
# or
dbmate down
```

### Run

```bash
go run cmd/api/main.go
```

## Usage Flow

1. **Authenticate with Google**
   ```bash
   # Visit in browser
   http://localhost:8080/api/auth/google/start
   # Complete OAuth, receive JWT token
   ```

2. **Register Sheets You Want to Access**
   ```bash
   # Extract Sheet ID from URL: https://docs.google.com/spreadsheets/d/SHEET_ID_HERE/edit
   
   curl -X POST http://localhost:8080/api/sheets/register \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
       "sheet_name": "Sales Data",
       "description": "Q1 2026 sales tracking"
     }'
   ```

3. **Access Registered Sheet Data**
   ```bash
   curl -X POST http://localhost:8080/api/sheets/data \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
       "range": "Sheet1!A1:Z100"
     }'
   ```

4. **List Your Registered Sheets**
   ```bash
   curl -X GET http://localhost:8080/api/sheets/registered \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

5. **Remove Sheet from Allowed List**
   ```bash
   curl -X DELETE http://localhost:8080/api/sheets/registered/SHEET_ID_HERE \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## Security

- **JWT tokens**: Used for user session management and authenticated operations
- **OAuth scope**: `spreadsheets.readonly` (read-only, clean consent screen)
- **Application-level access control**: Database tracks which sheets each user has registered
- **No password authentication**: Only Google OAuth
- **HTTPS recommended**: Use HTTPS in production

## Removed Features

This project has been simplified from the original template:
- ❌ No password authentication
- ❌ No refresh tokens
- ❌ No role-based access control
- ❌ No data caching (reads directly from Google Sheets)
