package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/congdv/go-auth/api/internal/config"
	"github.com/congdv/go-auth/api/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
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
			Scopes:       []string{"openid", "email", "profile", sheets.SpreadsheetsReadonlyScope},
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

	url := h.oauth.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"), oauth2.SetAuthURLParam("access_type", "offline"))

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user"})
		return
	}

	// Store Google OAuth tokens for future API access
	if err := h.auth.UpdateGoogleTokens(c.Request.Context(), user.ID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store oauth tokens"})
		return
	}

	// Update user object with the tokens
	user.GoogleAccessToken = &token.AccessToken
	user.GoogleRefreshToken = &token.RefreshToken
	user.GoogleTokenExpiry = &token.Expiry

	accessToken, exp, err := h.auth.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      accessToken,
		"expires_at": exp,
		"user":       user,
	})

	// when we have front end code
	// url := h.cfg.FrontendOrigin + "/oauth/callback?token=" + accessToken
	// c.Redirect(http.StatusFound, url)
}

func randomState(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
