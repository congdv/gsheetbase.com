-- migrate:up
-- =============================================================================
-- Add Authentication Type Support to Allowed Sheets
-- =============================================================================
-- This migration adds auth_type (none, bearer, basic) support to allowed_sheets table.
-- Allows sheets to be protected with Bearer tokens or Basic auth credentials.
-- =============================================================================

-- Create auth_type enum
CREATE TYPE auth_type AS ENUM ('none', 'bearer', 'basic');

-- Add auth_type and credential columns to allowed_sheets
ALTER TABLE allowed_sheets
  ADD COLUMN auth_type auth_type NOT NULL DEFAULT 'none',
  ADD COLUMN auth_bearer_token TEXT,
  ADD COLUMN auth_basic_username TEXT,
  ADD COLUMN auth_basic_password_hash TEXT;

-- Create indexes for efficient lookups by auth credentials
-- Bearer token lookup (when auth_type = 'bearer')
CREATE INDEX idx_allowed_sheets_bearer_token ON allowed_sheets(auth_bearer_token) 
  WHERE auth_type = 'bearer' AND auth_bearer_token IS NOT NULL;

-- Basic username lookup (when auth_type = 'basic')
CREATE INDEX idx_allowed_sheets_basic_username ON allowed_sheets(auth_basic_username) 
  WHERE auth_type = 'basic' AND auth_basic_username IS NOT NULL;

-- -- Add constraint: bearer auth requires bearer token
-- ALTER TABLE allowed_sheets
--   ADD CONSTRAINT check_bearer_token_required 
--     CHECK (auth_type != 'bearer' OR auth_bearer_token IS NOT NULL);

-- -- Add constraint: basic auth requires username and password hash
-- ALTER TABLE allowed_sheets
--   ADD CONSTRAINT check_basic_credentials_required 
--     CHECK (auth_type != 'basic' OR (auth_basic_username IS NOT NULL AND auth_basic_password_hash IS NOT NULL));

-- -- Add constraint: no bearer token if not bearer auth
-- ALTER TABLE allowed_sheets
--   ADD CONSTRAINT check_bearer_token_exclusive 
--     CHECK (auth_type = 'bearer' OR auth_bearer_token IS NULL);

-- -- Add constraint: no basic credentials if not basic auth
-- ALTER TABLE allowed_sheets
--   ADD CONSTRAINT check_basic_credentials_exclusive 
--     CHECK (auth_type = 'basic' OR (auth_basic_username IS NULL AND auth_basic_password_hash IS NULL));

-- migrate:down
-- =============================================================================
-- ROLLBACK
-- =============================================================================

-- Drop constraints
ALTER TABLE allowed_sheets
  DROP CONSTRAINT IF EXISTS check_bearer_token_exclusive,
  DROP CONSTRAINT IF EXISTS check_basic_credentials_exclusive,
  DROP CONSTRAINT IF EXISTS check_bearer_token_required,
  DROP CONSTRAINT IF EXISTS check_basic_credentials_required;

-- Drop indexes
DROP INDEX IF EXISTS idx_allowed_sheets_basic_username;
DROP INDEX IF EXISTS idx_allowed_sheets_bearer_token;

-- Drop columns
ALTER TABLE allowed_sheets
  DROP COLUMN IF EXISTS auth_basic_password_hash,
  DROP COLUMN IF EXISTS auth_basic_username,
  DROP COLUMN IF EXISTS auth_bearer_token,
  DROP COLUMN IF EXISTS auth_type;

-- Drop enum type
DROP TYPE IF EXISTS auth_type;
