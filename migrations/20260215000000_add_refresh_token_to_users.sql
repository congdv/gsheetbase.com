-- migrate:up
-- Add refresh token support to users table
-- Allows issuance of JWT refresh tokens for longer-lived sessions

ALTER TABLE users ADD COLUMN refresh_token_hash TEXT;
ALTER TABLE users ADD COLUMN refresh_token_expiry TIMESTAMPTZ;

-- Index for refresh token lookups (though we'll use hashing in practice)
CREATE INDEX idx_users_refresh_token_expiry ON users(refresh_token_expiry) WHERE refresh_token_hash IS NOT NULL;

-- migrate:down
DROP INDEX IF EXISTS idx_users_refresh_token_expiry;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token_expiry;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token_hash;
