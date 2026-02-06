-- migrate:up
-- Add scope tracking to users table
ALTER TABLE users ADD COLUMN google_scopes TEXT[];

-- Add write operation support to allowed_sheets table
ALTER TABLE allowed_sheets 
  ADD COLUMN allow_write BOOLEAN DEFAULT false,
  ADD COLUMN allowed_methods TEXT[] DEFAULT ARRAY['GET'];

-- Create index for efficient scope lookups
CREATE INDEX idx_users_scopes ON users USING GIN(google_scopes);

-- Update existing sheets to have GET method
UPDATE allowed_sheets SET allowed_methods = ARRAY['GET'] WHERE allowed_methods IS NULL;

-- migrate:down
-- Remove scope and write operation support
DROP INDEX IF EXISTS idx_users_scopes;
ALTER TABLE allowed_sheets DROP COLUMN IF EXISTS allowed_methods;
ALTER TABLE allowed_sheets DROP COLUMN IF EXISTS allow_write;
ALTER TABLE users DROP COLUMN IF EXISTS google_scopes;
