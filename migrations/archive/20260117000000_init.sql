-- migrate:up
-- Enable useful extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- Users (simplified for OAuth only)
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email CITEXT UNIQUE NOT NULL,
  provider TEXT NOT NULL,
  provider_id TEXT NOT NULL,
  google_access_token TEXT,
  google_refresh_token TEXT,
  google_token_expiry TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (provider, provider_id)
);

-- Useful indexes
CREATE INDEX idx_users_provider_id ON users(provider, provider_id);

-- migrate:down
DROP TABLE IF EXISTS users;
