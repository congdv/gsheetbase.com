#!/bin/bash

# Test subscription and quota enforcement
# Usage: ./scripts/test_subscription.sh [API_KEY]

set -e

API_KEY="${1:-YOUR_API_KEY}"
WORKER_URL="${WORKER_URL:-http://localhost:8081}"
WEB_URL="${WEB_URL:-http://localhost:8080}"

echo "🧪 Testing Subscription & Quota System"
echo "======================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test 1: Get available plans (public)
echo -e "${YELLOW}Test 1: Get Available Plans (Public)${NC}"
curl -s "${WEB_URL}/api/subscription/plans" | jq '.' || echo "Failed"
echo ""

# Test 2: Make a GET request and check quota headers
echo -e "${YELLOW}Test 2: GET Request - Check Quota Headers${NC}"
RESPONSE=$(curl -s -i "${WORKER_URL}/v1/${API_KEY}")
echo "$RESPONSE" | grep -E "X-RateLimit-|X-Daily-Quota-|X-Monthly-Quota-"
echo ""

# Test 3: Test rate limiting (send multiple requests quickly)
echo -e "${YELLOW}Test 3: Rate Limit Test (GET)${NC}"
echo "Sending 25 GET requests rapidly..."
SUCCESS_COUNT=0
RATE_LIMITED_COUNT=0

for i in {1..25}; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${WORKER_URL}/v1/${API_KEY}")
  
  if [ "$STATUS" = "200" ] || [ "$STATUS" = "304" ]; then
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    echo -n "."
  elif [ "$STATUS" = "429" ]; then
    RATE_LIMITED_COUNT=$((RATE_LIMITED_COUNT + 1))
    echo -n "X"
  else
    echo -n "?"
  fi
done

echo ""
echo -e "${GREEN}Successful: ${SUCCESS_COUNT}${NC}"
echo -e "${RED}Rate Limited: ${RATE_LIMITED_COUNT}${NC}"
echo ""

# Test 4: Test UPDATE rate limit (slower than GET)
echo -e "${YELLOW}Test 4: Rate Limit Test (UPDATE)${NC}"
echo "Sending 10 POST requests rapidly..."
UPDATE_SUCCESS=0
UPDATE_LIMITED=0

for i in {1..10}; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${WORKER_URL}/v1/${API_KEY}" \
    -H "Content-Type: application/json" \
    -d '{"data": [{"test": "value"}]}')
  
  if [ "$STATUS" = "200" ] || [ "$STATUS" = "201" ]; then
    UPDATE_SUCCESS=$((UPDATE_SUCCESS + 1))
    echo -n "."
  elif [ "$STATUS" = "429" ]; then
    UPDATE_LIMITED=$((UPDATE_LIMITED + 1))
    echo -n "X"
  else
    echo -n "?"
  fi
  
  sleep 0.1 # Small delay to avoid immediate rate limiting
done

echo ""
echo -e "${GREEN}Successful: ${UPDATE_SUCCESS}${NC}"
echo -e "${RED}Rate Limited: ${UPDATE_LIMITED}${NC}"
echo ""

# Test 5: Check rate limit headers after hitting limit
echo -e "${YELLOW}Test 5: Rate Limit Response Details${NC}"
curl -s -i "${WORKER_URL}/v1/${API_KEY}" | head -n 20
echo ""

# Summary
echo "======================================="
echo -e "${GREEN}✅ Subscription system test complete${NC}"
echo ""
echo "📊 Next Steps:"
echo "  1. Check your current plan: GET ${WEB_URL}/api/subscription/plan"
echo "  2. Check usage stats: GET ${WEB_URL}/api/subscription/usage"
echo "  3. Run migration: psql \$DATABASE_URL < migrations/20260124000001_add_subscription_plans.sql"
echo ""
