CREATE TABLE plans (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_name             TEXT NOT NULL,  -- 'free', 'starter', 'pro', 'enterprise'
    
    -- Rate Limits (per minute)
    get_rate_limit        INTEGER NOT NULL,
    update_rate_limit     INTEGER NOT NULL,
    
    -- Quotas
    daily_update_quota    INTEGER NOT NULL,
    monthly_get_quota     INTEGER NOT NULL,
    monthly_update_quota  INTEGER NOT NULL,
    
    -- Pricing (in cents USD)
    monthly_price_cents   INTEGER NOT NULL,
    annual_price_cents    INTEGER NOT NULL,
    
    -- Features
    cache_min_ttl         INTEGER NOT NULL,  -- seconds
    custom_domain         BOOLEAN NOT NULL DEFAULT false,
    priority_support      BOOLEAN NOT NULL DEFAULT false,
    
    -- Temporal validity
    effective_from        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    effective_to          TIMESTAMPTZ,  -- NULL means currently active
    
    -- Audit
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT plans_name_effective_from_unique UNIQUE (plan_name, effective_from),
    CONSTRAINT plans_effective_dates_check CHECK (effective_to IS NULL OR effective_to > effective_from)
);

CREATE INDEX idx_plans_current ON plans (plan_name, effective_from, effective_to) 
    WHERE effective_to IS NULL;

CREATE INDEX idx_plans_temporal ON plans (plan_name, effective_from, effective_to);