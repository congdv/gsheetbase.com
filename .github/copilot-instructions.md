# Copilot Instructions – GSheetBase

Act as a **Senior Full-Stack Engineer and Software Architect**.
You are helping build **GSheetBase**, a SaaS platform that converts **Google Sheets into REST APIs**.

The goal of these instructions is to guide GitHub Copilot (and any AI pair programmer) to generate **correct, simple, production-ready code** that aligns with the architecture and product vision.

---

## 🧠 Product Overview

**GSheetBase** allows users to:

* Sign in with Google (OAuth)
* Select specific Google Sheets (not all files)
* Generate public or private REST API endpoints
* Fetch Google Sheet data as clean, structured JSON
* Control access, caching, and performance

---

## 🏗️ Architecture Overview

### Monorepo Structure (Railway-friendly)

```
/gsheetbase
├── api/              # Main backend (Gin)
│   ├── cmd/
│   ├── internal/
│   │   ├── auth/     # Google OAuth, JWT, session handling
│   │   ├── users/    # User & project metadata
│   │   ├── projects/ # Sheet → API configuration
│   │   ├── http/     # HTTP handlers (Gin controllers)
│   │   └── db/       # PostgreSQL access
│   └── main.go
│
├── worker/           # Sheet API worker service
│   ├── internal/
│   │   ├── sheets/   # Google Sheets API access
│   │   ├── cache/    # Redis caching
│   │   └── fetcher/  # Fetch + normalize sheet data
│   └── main.go
│
├── ui/               # Frontend (React)
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/ # Axios API clients
│   │   └── styles/
│   └── vite.config.ts
│
└── README.md
```

---

## ⚙️ Backend Guidelines

### Tech Stack

* Language: **Go (Golang)**
* Web Framework: **Gin**
* Pattern: **Light BFF (Backend For Frontend)**
* Database: **PostgreSQL** (user, project, metadata)
* Cache: **Redis** (API responses)

### Backend Services

#### 1. Main Backend (`/api`)

Responsible for:

* Authentication (Google OAuth)
* User accounts & projects
* API key management
* Access control (public/private APIs)
* Serving configuration to workers

#### 2. Worker Service (`/worker`)

Responsible for:

* Fetching Google Sheets data
* Using stored OAuth tokens (scoped)
* Normalizing sheet rows → JSON
* Caching responses in Redis
* Handling high-volume API traffic

Workers must be **stateless** and horizontally scalable.

---

## 🔐 Google OAuth Rules

* Use **OAuth 2.0** with incremental authorization
* Default scope:

  ```
  https://www.googleapis.com/auth/spreadsheets.readonly
  ```
* App **must not** access all Google Drive files
* User explicitly selects which spreadsheet to connect
* Store:

  * Google `spreadsheetId`
  * Access token (encrypted)
  * Refresh token

Never assume access to all user spreadsheets.

---

## 🌐 API Design Rules

* RESTful
* Versioned (`/v1`)
* JSON-only responses
* Clear error messages

Example response:

```json
{
  "data": [
    { "name": "Apple", "price": 1.2 },
    { "name": "Banana", "price": 0.8 }
  ],
  "meta": {
    "rows": 2,
    "cached": true
  }
}
```

---

## 🎨 Frontend Guidelines

### Tech Stack

* React + Vite
* TypeScript
* **Ant Design (antd)**
* **TanStack Query** (data fetching & caching)
* Axios (HTTP client)

### UI Principles

* **Mobile-first** responsive design
* Use Ant Design grid and breakpoints
* Clean SaaS-style dashboard UI
* Avoid over-engineering

### Frontend Responsibilities

* Google OAuth flow
* Project & API management UI
* Display generated API endpoints
* Show live preview of JSON data

---

## 📡 Frontend Data Rules

* Use **TanStack Query** for all server state
* Axios instances must:

  * Handle auth headers
  * Handle 401 / 403 globally
* No direct Google API calls from frontend

---

## 🧪 Coding Instructions for Copilot

When generating code:

* ✅ Provide the **simplest correct solution first**
* ✅ Follow existing folder structure
* ✅ Explain trade-offs briefly if needed
* ❌ Do not rewrite unrelated code
* ❌ Do not introduce unnecessary abstractions
* ❌ Do not assume the app is running locally

Prefer:

* Clear function names
* Explicit types (especially in Go & TS)
* Readability over cleverness

---

## 🚀 Deployment Assumptions

* Deployed on **Railway**
* Each folder (`/api`, `/worker`, `/ui`) is a separate service
* UI served as static assets or standalone frontend service
* APIs exposed via subdomains:

  * `api.gsheetbase.com`
  * `app.gsheetbase.com`

---

## 🎯 Final Goal

Build a **scalable, developer-friendly platform** where:

* Non-technical users can create APIs from Google Sheets
* Developers can rely on stable, fast, cached endpoints
* The system scales horizontally with minimal complexity

Always optimize for **clarity, security, and maintainability**.
