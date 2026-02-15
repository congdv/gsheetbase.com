package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/http/middleware"
	"gsheetbase/web/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthHandler struct {
	auth  services.AuthService
	cfg   *config.Config
	oauth *oauth2.Config
}

func NewGoogleAuthHandler(auth services.AuthService, cfg *config.Config) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		auth: auth,
		cfg:  cfg,
		oauth: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectUrl,
			Scopes:       []string{"email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *GoogleAuthHandler) Start(c *gin.Context) {
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "google auth not configured"})
		return
	}

	state := randomState(32)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, int((5 * time.Minute).Seconds()), "/api/auth/google", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	// Use include_granted_scopes to preserve previously granted permissions
	url := h.oauth.AuthCodeURL(state,
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)

	c.Redirect(http.StatusFound, url)
}

func (h *GoogleAuthHandler) Callback(c *gin.Context) {
	stateParam := c.Query("state")
	code := c.Query("code")
	stateCookie, err := c.Cookie("oauth_state")

	if err != nil || stateCookie == "" || stateParam == "" || stateCookie != stateParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	ctx := context.Background()
	token, err := h.oauth.Exchange(ctx, code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "oauth exchange failed"})
		return
	}

	client := h.oauth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch user info"})
		return
	}

	defer resp.Body.Close()

	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid user info"})
		return
	}

	if info.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not provided by provider"})
		return
	}

	user, err := h.auth.FindOrCreateOauthUser(c.Request.Context(), info.Email, "google", info.Sub)
	if err != nil {
		fmt.Printf("ERROR FindOrCreateOauthUser: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user", "details": err.Error()})
		return
	}

	// Store Google OAuth tokens for future API access
	if err := h.auth.UpdateGoogleTokens(c.Request.Context(), user.ID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		fmt.Printf("ERROR UpdateGoogleTokens: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store oauth tokens", "details": err.Error()})
		return
	}

	// Extract granted scopes from the scope query parameter
	// Gin automatically decodes the query parameter, so spaces will be actual spaces
	scopeParam := c.Query("scope")
	var scopes []string
	if scopeParam != "" {
		// Scopes are space-separated in the query param (already decoded by Gin)
		scopes = strings.Split(scopeParam, " ")
	} else {
		// Fallback to default scopes
		scopes = []string{"openid", "email", "profile"}
	}

	fmt.Printf("DEBUG: Storing scopes: %v\n", scopes)
	if err := h.auth.UpdateGoogleScopes(c.Request.Context(), user.ID, scopes); err != nil {
		fmt.Printf("ERROR UpdateGoogleScopes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store oauth scopes", "details": err.Error()})
		return
	}

	// Update user object with the tokens
	user.GoogleAccessToken = &token.AccessToken
	user.GoogleRefreshToken = &token.RefreshToken
	user.GoogleTokenExpiry = &token.Expiry
	user.GoogleScopes = scopes

	accessToken, exp, err := h.auth.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue access token"})
		return
	}

	// Set access token as HTTP-only cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, int(time.Until(exp).Seconds()), "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	// Generate and set refresh token cookie
	refreshToken, err := h.auth.GenerateRefreshToken(c.Request.Context(), user.ID)
	if err != nil {
		// Log but don't fail - refresh token is optional for now
		fmt.Printf("ERROR GenerateRefreshToken: %v\n", err)
	} else {
		refreshTokenExpiry := time.Now().AddDate(0, 0, h.cfg.JWTRefreshTTLDays)
		c.SetCookie("refresh_token", refreshToken, int(time.Until(refreshTokenExpiry).Seconds()), "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	}

	// Check if this is incremental auth (from popup)
	incrementalAuth, _ := c.Cookie("incremental_auth")
	if incrementalAuth == "true" {
		// Clear the incremental auth cookie
		c.SetCookie("incremental_auth", "", -1, "/api/auth/google", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

		// Close the popup window with a success message
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			<head><title>Authorization Complete</title></head>
			<body>
				<script>
					window.close();
				</script>
				<p>Authorization successful! You can close this window.</p>
			</body>
			</html>
		`)
		return
	}

	// Regular OAuth flow - redirect to frontend callback
	url := h.cfg.FrontendOrigin + "/oauth/callback"
	c.Redirect(http.StatusFound, url)
}

func randomState(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (h *GoogleAuthHandler) RequestAdditionalScopes(c *gin.Context) {
	var req struct {
		Scopes []string `json:"scopes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scopes required"})
		return
	}

	// Get current user ID from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.auth.Me(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Merge existing scopes with new requested scopes
	scopeMap := make(map[string]bool)
	for _, s := range user.GoogleScopes {
		scopeMap[s] = true
	}
	for _, s := range req.Scopes {
		scopeMap[s] = true
	}

	allScopes := make([]string, 0, len(scopeMap))
	for s := range scopeMap {
		allScopes = append(allScopes, s)
	}

	// Create OAuth config with merged scopes
	oauth := &oauth2.Config{
		ClientID:     h.cfg.GoogleClientID,
		ClientSecret: h.cfg.GoogleClientSecret,
		RedirectURL:  h.cfg.GoogleRedirectUrl,
		Scopes:       allScopes,
		Endpoint:     google.Endpoint,
	}

	state := randomState(32)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, int((5 * time.Minute).Seconds()), "/api/auth/google", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie("incremental_auth", "true", int((5 * time.Minute).Seconds()), "/api/auth/google", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	// Force consent to show new permissions, include previously granted scopes
	url := oauth.AuthCodeURL(state,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)

	c.JSON(http.StatusOK, gin.H{"auth_url": url})
}
