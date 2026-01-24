# Test rate limiting (should get 429 after 60 requests)
for i in {1..65}; do 
  echo "Request $i:"
  curl -i http://localhost:8081/v1/gsheet_TX5MEZq-gkK0R2YpvLYxOHUSKOOOR1Yu | grep -E "HTTP|X-RateLimit"
done