# GSheetBase

A platform that allows users to connect their Google Sheets via OAuth and generate REST API endpoints that return sheet data as clean JSON.

## Features

- ✅ **Google OAuth Authentication** - Simple, secure sign-in with Google
- ✅ **Sheet Registration** - Whitelist specific sheets for API access
- ✅ **Read-Only Access** - Uses `spreadsheets.readonly` OAuth scope
- ✅ **JWT Sessions** - Secure token-based authentication
- ✅ **No Data Storage** - Direct Google Sheets API access (sheets data not stored)
- ✅ **Clean OAuth Consent** - No scary "edit/delete" warnings

## Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: PostgreSQL (user metadata only)
- **Authentication**: Google OAuth 2.0 + JWT

### Frontend
- **Framework**: React with Vite (TypeScript)
- **UI Library**: Ant Design
- **Styling**: Styled Components + Tailwind CSS
- **State Management**: TanStack Query

## Project Structure

```
/web
├── /cmd/api          # Backend entry point
├── /internal         # Backend business logic
│   ├── /config       # Configuration
│   ├── /database     # Database connection
│   ├── /http         # HTTP handlers & middleware
│   ├── /models       # Data models
│   ├── /repository   # Database queries
│   └── /services     # Business logic
└── /ui               # React frontend
    └── /src
        ├── /components
        ├── /context     # Auth context
        ├── /lib         # Axios, React Query
        └── /pages       # Login, Dashboard
```

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Google OAuth credentials

### 1. Setup Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable Google Sheets API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/api/auth/google/callback`
6. Download credentials or copy Client ID and Secret

### 2. Setup Environment Variables

Create `.env` in the project root:

```env
# Server
PORT=8080

# Database
DATABASE_URL=postgres://user:password@localhost:5432/gsheetbase?sslmode=disable

# JWT
JWT_ACCESS_SECRET=your-random-secret-key-change-this
JWT_ACCESS_TTL_MINUTES=60

# CORS
FRONTEND_ORIGIN=http://localhost:5173
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

Create `web/ui/.env`:

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

### 3. Database Setup

```bash
# Run migrations
make migrate-up
```

### 4. Run Backend

```bash
# Install dependencies
go mod download

# Run server
make run

# Or directly
go run web/cmd/api/main.go
```

Backend will run on `http://localhost:8080`

### 5. Run Frontend

```bash
cd web/ui

# Install dependencies
npm install

# Run dev server
npm run dev
```

Frontend will run on `http://localhost:5173`

## How to Use

### 1. Sign In with Google

- Navigate to `http://localhost:5173`
- Click "Continue with Google"
- Grant permissions to view your Google Sheets
- You'll be redirected to the dashboard

### 2. Register a Sheet

1. Click "Register Sheet" button
2. Paste your Google Sheets URL or Sheet ID
   - Example URL: `https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit`
   - Sheet ID: `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms`
3. Optionally add a name and description
4. Click OK

### 3. Access Sheet Data via API

Once registered, you can access sheet data:

```bash
# Get sheet data (requires authentication cookie)
curl -X POST http://localhost:8080/api/sheets/data \
  -H "Content-Type: application/json" \
  -d '{
    "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "range": "Sheet1!A1:E10"
  }'
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/google/start` | Start Google OAuth flow |
| GET | `/api/auth/google/callback` | OAuth callback (redirects to frontend) |
| GET | `/api/auth/me` | Get current user info (requires auth) |
| POST | `/api/auth/logout` | Clear session cookie |

### Sheet Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/sheets/register` | Register a new sheet |
| GET | `/api/sheets/registered` | List registered sheets |
| DELETE | `/api/sheets/registered/:sheet_id` | Remove a registered sheet |
| POST | `/api/sheets/data` | Get data from a registered sheet |

### Example: Get Sheet Data

```bash
POST /api/sheets/data
Content-Type: application/json

{
  "sheet_id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
  "range": "Sheet1!A1:E10"
}
```

Response:
```json
{
  "data": [
    ["Name", "Email", "Status", "Date", "Amount"],
    ["John Doe", "john@example.com", "Active", "2026-01-15", "100"],
    ["Jane Smith", "jane@example.com", "Pending", "2026-01-16", "200"]
  ]
}
```

## Security Model

### Application-Level Access Control

This app uses **application-level per-sheet access control** with the `spreadsheets.readonly` OAuth scope:

1. ✅ Users authenticate with Google OAuth
2. ✅ Users explicitly register specific sheets they want to access
3. ✅ Only registered sheets can be accessed via the API
4. ✅ Database tracks which sheets each user has authorized
5. ✅ Read-only access enforced by OAuth scope

### Why This Approach?

- **Clean OAuth Consent**: Users see "View your Google Spreadsheets" instead of scary "edit and delete" warnings
- **Explicit Control**: Users choose exactly which sheets to expose
- **Audit Trail**: Database records all registered sheets
- **Security**: Can't accidentally access or modify wrong sheets

## Development

### Available Make Commands

```bash
make run          # Run the backend server
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
make build        # Build the application
```

### Frontend Scripts

```bash
npm run dev       # Start dev server
npm run build     # Build for production
npm run preview   # Preview production build
```

## Deployment

### Backend (Railway/Render/Fly.io)

1. Set environment variables
2. Run migrations
3. Deploy Go application

### Frontend (Vercel/Netlify)

1. Set `VITE_API_BASE_URL` to your production API URL
2. Build and deploy

## Future Enhancements

- [ ] API key generation for programmatic access
- [ ] Rate limiting per user
- [ ] Webhook support for sheet changes
- [ ] Public/private sheet endpoints
- [ ] Caching with Redis
- [ ] Sheet data transformation (JSON, CSV, etc.)

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
