-- migrate:up

-- Add subscription plan enum
CREATE TYPE subscription_plan AS ENUM ('free', 'starter', 'pro', 'enterprise');

-- Add subscription billing period enum
CREATE TYPE billing_period AS ENUM ('monthly', 'annual');

-- Add subscription fields to users table
ALTER TABLE users ADD COLUMN subscription_plan subscription_plan DEFAULT 'free' NOT NULL;
ALTER TABLE users ADD COLUMN billing_period billing_period DEFAULT 'monthly';
ALTER TABLE users ADD COLUMN subscription_status TEXT DEFAULT 'active';
ALTER TABLE users ADD COLUMN subscription_started_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN subscription_ends_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN stripe_customer_id TEXT;
ALTER TABLE users ADD COLUMN stripe_subscription_id TEXT;

-- Add indexes for billing queries
CREATE INDEX idx_users_subscription_plan ON users(subscription_plan);
CREATE INDEX idx_users_subscription_status ON users(subscription_status);
CREATE INDEX idx_users_stripe_customer_id ON users(stripe_customer_id);

-- Add index for monthly usage rollups (using date range queries)
CREATE INDEX idx_api_usage_daily_monthly ON api_usage_daily(user_id, request_date, method);

-- migrate:down
DROP INDEX IF EXISTS idx_api_usage_daily_monthly;
DROP INDEX IF EXISTS idx_users_stripe_customer_id;
DROP INDEX IF EXISTS idx_users_subscription_status;
DROP INDEX IF EXISTS idx_users_subscription_plan;

ALTER TABLE users DROP COLUMN IF EXISTS stripe_subscription_id;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_customer_id;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_ends_at;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_started_at;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_status;
ALTER TABLE users DROP COLUMN IF EXISTS billing_period;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_plan;

DROP TYPE IF EXISTS billing_period;
DROP TYPE IF EXISTS subscription_plan;
