-- migrate:up
-- Add API customization columns to allowed_sheets table
ALTER TABLE allowed_sheets
    ADD COLUMN api_key TEXT UNIQUE,
    ADD COLUMN is_public BOOLEAN DEFAULT false NOT NULL,
    ADD COLUMN default_range TEXT DEFAULT 'Sheet1',
    ADD COLUMN use_first_row_as_header BOOLEAN DEFAULT true NOT NULL;

-- Add index for API key lookups
CREATE INDEX idx_allowed_sheets_api_key ON allowed_sheets(api_key) WHERE api_key IS NOT NULL;

-- Add index for public sheets
CREATE INDEX idx_allowed_sheets_public ON allowed_sheets(is_public) WHERE is_public = true;

-- migrate:down
DROP INDEX IF EXISTS idx_allowed_sheets_public;
DROP INDEX IF EXISTS idx_allowed_sheets_api_key;
ALTER TABLE allowed_sheets 
    DROP COLUMN IF EXISTS use_first_row_as_header,
    DROP COLUMN IF EXISTS default_range,
    DROP COLUMN IF EXISTS is_public,
    DROP COLUMN IF EXISTS api_key;
