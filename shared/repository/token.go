package repository

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateBearerToken returns a cryptographically secure bearer token
// prefixed with "gsheet_". It uses URL-safe base64 encoding and does
// not assume a fixed encoded length to avoid slicing panics.
func GenerateBearerToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "gskey_fallback_token"
	}
	s := base64.URLEncoding.EncodeToString(b)
	return "gskey_" + s
}
