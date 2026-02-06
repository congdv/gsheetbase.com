-- migrate:up
-- =============================================================================
-- Gsheetbase Initial Schema - First Release
-- =============================================================================
-- This migration consolidates all schema changes for the first release.
-- It includes: users with OAuth, allowed sheets, usage tracking, and subscriptions.
-- Removed: rate_limit_override (unused in codebase)
-- =============================================================================

-- Enable useful extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- =============================================================================
-- ENUMS
-- =============================================================================

-- Subscription plan tiers
CREATE TYPE subscription_plan AS ENUM ('free', 'starter', 'pro', 'enterprise');

-- Billing period options
CREATE TYPE billing_period AS ENUM ('monthly', 'annual');

-- =============================================================================
-- USERS TABLE
-- =============================================================================
-- Stores user accounts with Google OAuth and subscription information

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email CITEXT UNIQUE NOT NULL,
  provider TEXT NOT NULL,
  provider_id TEXT NOT NULL,
  
  -- Google OAuth tokens
  google_access_token TEXT,
  google_refresh_token TEXT,
  google_token_expiry TIMESTAMPTZ,
  google_scopes TEXT[],
  
  -- Subscription & Billing
  subscription_plan subscription_plan DEFAULT 'free' NOT NULL,
  billing_period billing_period DEFAULT 'monthly',
  subscription_status TEXT DEFAULT 'active',
  subscription_started_at TIMESTAMPTZ,
  subscription_ends_at TIMESTAMPTZ,
  stripe_customer_id TEXT,
  stripe_subscription_id TEXT,
  
  -- Timestamps
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  UNIQUE (provider, provider_id)
);

-- Users indexes
CREATE INDEX idx_users_provider_id ON users(provider, provider_id);
CREATE INDEX idx_users_scopes ON users USING GIN(google_scopes);
CREATE INDEX idx_users_subscription_plan ON users(subscription_plan);
CREATE INDEX idx_users_subscription_status ON users(subscription_status);
CREATE INDEX idx_users_stripe_customer_id ON users(stripe_customer_id);

-- =============================================================================
-- ALLOWED SHEETS TABLE
-- =============================================================================
-- Stores user's registered Google Sheets with API configuration

CREATE TABLE allowed_sheets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  sheet_id TEXT NOT NULL,
  sheet_name TEXT,
  description TEXT,
  
  -- API Configuration
  api_key TEXT UNIQUE,
  is_public BOOLEAN DEFAULT false NOT NULL,
  default_range TEXT DEFAULT 'Sheet1',
  use_first_row_as_header BOOLEAN DEFAULT true NOT NULL,
  
  -- Write Access Control
  allow_write BOOLEAN DEFAULT false,
  allowed_methods TEXT[] DEFAULT ARRAY['GET'],
  
  -- Timestamps
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  UNIQUE (user_id, sheet_id)
);

-- Allowed sheets indexes
CREATE INDEX idx_allowed_sheets_user_id ON allowed_sheets(user_id);
CREATE INDEX idx_allowed_sheets_sheet_id ON allowed_sheets(sheet_id);
CREATE INDEX idx_allowed_sheets_api_key ON allowed_sheets(api_key) WHERE api_key IS NOT NULL;
CREATE INDEX idx_allowed_sheets_public ON allowed_sheets(is_public) WHERE is_public = true;

-- =============================================================================
-- API USAGE DAILY TABLE
-- =============================================================================
-- Tracks API usage per day for quota enforcement and analytics

CREATE TABLE api_usage_daily (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_key TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sheet_id UUID NOT NULL REFERENCES allowed_sheets(id) ON DELETE CASCADE,
    request_date DATE NOT NULL,
    method TEXT NOT NULL, -- GET, POST, PUT, PATCH
    request_count INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Ensure uniqueness per API key, date, and method
    UNIQUE (api_key, request_date, method)
);

-- API usage indexes
CREATE INDEX idx_api_usage_daily_api_key ON api_usage_daily(api_key);
CREATE INDEX idx_api_usage_daily_user_id ON api_usage_daily(user_id);
CREATE INDEX idx_api_usage_daily_sheet_id ON api_usage_daily(sheet_id);
CREATE INDEX idx_api_usage_daily_date ON api_usage_daily(request_date DESC);
CREATE INDEX idx_api_usage_daily_composite ON api_usage_daily(api_key, request_date DESC);
CREATE INDEX idx_api_usage_daily_monthly ON api_usage_daily(user_id, request_date, method);

-- =============================================================================
-- COMMENTS
-- =============================================================================

COMMENT ON TABLE users IS 'User accounts with Google OAuth and subscription management';
COMMENT ON TABLE allowed_sheets IS 'Google Sheets registered as APIs with access configuration';
COMMENT ON TABLE api_usage_daily IS 'Daily API request counts for quota enforcement and analytics';

COMMENT ON COLUMN users.google_scopes IS 'Array of Google API scopes granted by user';
COMMENT ON COLUMN allowed_sheets.allowed_methods IS 'HTTP methods permitted for this API (GET, POST, PUT, PATCH)';
COMMENT ON COLUMN api_usage_daily.request_count IS 'Total requests for this API key, date, and method';

-- migrate:down
-- =============================================================================
-- ROLLBACK
-- =============================================================================

DROP TABLE IF EXISTS api_usage_daily;
DROP TABLE IF EXISTS allowed_sheets;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS billing_period;
DROP TYPE IF EXISTS subscription_plan;

DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pgcrypto;
DROP EXTENSION IF EXISTS "uuid-ossp";
