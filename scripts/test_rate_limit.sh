# Test rate limiting (should get 429 after 60 requests)
for i in {1..65}; do 
  echo "Request $i:"
  curl -i http://localhost:8081/v1/gsheet_V-yRe3W_AeDG3O95v4Lm5skxBkKej6ic | grep -E "HTTP|X-RateLimit"
done