-- migrate:up
-- Allowed sheets table for application-level access control
CREATE TABLE allowed_sheets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  sheet_id TEXT NOT NULL,
  sheet_name TEXT,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, sheet_id)
);

CREATE INDEX idx_allowed_sheets_user_id ON allowed_sheets(user_id);
CREATE INDEX idx_allowed_sheets_sheet_id ON allowed_sheets(sheet_id);

-- migrate:down
DROP TABLE IF EXISTS allowed_sheets;
