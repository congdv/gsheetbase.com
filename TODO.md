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


# GraphQL-style REST API 

Example dataset:
| Asset             | Network access | Owner                        | Location             |  Notes                                                                 | Sensitivity     |
|-------------------|----------------|------------------------------|----------------------|-----------------------------------------------------------------------|-----------------|
| Patched           | Value          | Internet service provider (ISP) | On-premises          | Has a 2.4 GHz and 5 GHz connection. All devices on the home network connect to the 5 GHz frequency. | Confidential    |
| Desktop           | Occasional     | Homeowner                    | On-premises          | Contains private information, like photos.                            | Restricted      |
| Guest smartphone  | Occasional     | Friend                       | On and Off-premises  | Connects to my home network.                                          | Internal-only   |
| Printer           | Occasional     | Homeowner                    | On-premises          | Connects to my home network.                                          | Internal-only   |
| Smart TV          | Occasional     | Homeowner                    | On-premises          | Connects to my home network.                                          | Internal-only   |
| Google Home       | Continuous     | Homeowner                    | On-premises          | Connects to my home network.                                          | Internal-only   |


GET /v1/<api-key>?collection=Sheet1&fields=asset,location&where={"owner":"Homeowner"}&orderBy=asset

Request Parameters:

- collection: (string) The sheet name (e.g., assets).

- fields: (string) Optional, Comma-separated keys (e.g., asset,notes).

- limit: (number) Optional. Default 100.

- offset: (number) Optional. Default 0.

- orderBy: (string) Optional. e.g asset

- where: (json) search/filter. e.g {"owner":"Homeowner"}

Response:

with pagination
```
{
  "data": [
    { "asset": "Internet service provider", "owner": "Patched" },
    { "asset": "Desktop", "owner": "Homeowner" }
  ],
  "pagination": {
    "total": 6,
    "limit": 2,
    "offset": 0,
    "nextOffset": 2
  }
}
```

without pagination
```
{
  "data": [
    { "asset": "Internet service provider", "owner": "Patched" },
    { "asset": "Desktop", "owner": "Homeowner" }
  ]
}
```

POST /v1/<api-key>

Adding a new row to google sheet

Request body:

e.g:

with returing fields
```
{
  "collection": "assets",
  "data": {
    "asset": "Scanner",
    "owner": "Homeowner",
    "sensitivity": "Internal-only"
  },
  "returning": ["asset", "owner"]
}
```

without returning fields
```
{
  "collection": "assets",
  "data": {
    "asset": "Scanner",
    "owner": "Homeowner",
    "sensitivity": "Internal-only"
  }
}
```

Response:


201 with returning fields or all created fields if no returning fields


PUT /v1/<api-key>

it is like full replacement for the whole object, PATCH is partial updates

Request body

with returing fields

{
  "collection": "assets",
  "where": {"owner":"Homeowner"},
  "data": {
    "asset": "Desktop",
    "networkAccess": "Wired Only",
    "owner": "Homeowner",
    "location": "Office",
    "notes": "Completely reset via API",
    "sensitivity": "Restricted"
  }
  "returning": ["asset", "notes", "sensitivity"]
}

without returing fields

{
  "collection": "assets",
  "where": {"owner":"Homeowner"},
  "data": {
    "asset": "Desktop",
    "networkAccess": "Wired Only",
    "owner": "Homeowner",
    "location": "Office",
    "notes": "Completely reset via API",
    "sensitivity": "Restricted"
  }
}

Response:

Status: 200 OK if you return the updated resource


PATCH /v1/<api-key>

Request body

with returing fields

{
  "collection": "assets",
  "where": {"owner":"Homeowner"},
  "data": {
    "notes": "Updated memory to 32GB",
    "sensitivity": "Confidential"
  },
  "returning": ["asset", "notes", "sensitivity"]
}

without returing fields

{
  "collection": "assets",
  "where": {"owner":"Homeowner"},
  "data": {
    "notes": "Updated memory to 32GB",
    "sensitivity": "Confidential"
  },
}

Response:

Status: 200 OK if you return the updated resource


DELETE /v1/<api-key>?collection=assets&where={"owner":"Homeowner"}

Request parameter:
- collection: (string) The sheet name (e.g., assets).

- where: (json) search/filter. e.g {"owner":"Homeowner"}

Response


HTTP/1.1 204 No Content