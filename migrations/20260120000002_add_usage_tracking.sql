-- migrate:up
-- Add table for tracking API usage by date
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

-- Indexes for efficient querying
CREATE INDEX idx_api_usage_daily_api_key ON api_usage_daily(api_key);
CREATE INDEX idx_api_usage_daily_user_id ON api_usage_daily(user_id);
CREATE INDEX idx_api_usage_daily_sheet_id ON api_usage_daily(sheet_id);
CREATE INDEX idx_api_usage_daily_date ON api_usage_daily(request_date DESC);
CREATE INDEX idx_api_usage_daily_composite ON api_usage_daily(api_key, request_date DESC);

-- Add rate limit configuration to users table (optional: for tiered limits)
ALTER TABLE users ADD COLUMN rate_limit_per_minute INT DEFAULT 60 NOT NULL;
ALTER TABLE users ADD COLUMN rate_limit_burst INT DEFAULT 100 NOT NULL;

-- Add optional per-sheet rate limit override
ALTER TABLE allowed_sheets ADD COLUMN rate_limit_override INT;

-- migrate:down
DROP TABLE IF EXISTS api_usage_daily;
ALTER TABLE allowed_sheets DROP COLUMN IF EXISTS rate_limit_override;
ALTER TABLE users DROP COLUMN IF EXISTS rate_limit_burst;
ALTER TABLE users DROP COLUMN IF EXISTS rate_limit_per_minute;