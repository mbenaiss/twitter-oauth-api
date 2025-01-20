package twitter

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mbenaiss/twitter-oauth-api/models"
)

const (
	tokenURL    = "https://api.x.com/2/oauth2/token"
	authURL     = "https://x.com/i/oauth2/authorize"
	oauthScopes = "tweet.read users.read follows.read offline.access"
)

type Client struct {
	clientID     string
	clientSecret string
	codeVerifier string
	httpClient   *http.Client
	redirectURI  string
}

// NewClient creates a new OAuth client
func NewClient(clientID, clientSecret, redirectURI string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		redirectURI: redirectURI,
	}
}

func (c *Client) GetAuthURL() (string, error) {
	// Generate and store PKCE verifier
	verifier, err := generateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("error generating verifier: %v", err)
	}
	c.codeVerifier = verifier
	challenge := generateCodeChallenge(verifier)

	// Generate and store state
	state, err := generateState()
	if err != nil {
		return "", fmt.Errorf("error generating state: %v", err)
	}

	// Build authorization URL
	return fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
		authURL,
		c.clientID,
		c.redirectURI,
		oauthScopes,
		state,
		challenge,
	), nil
}

// ExchangeCodeForToken exchanges auth code for access token
func (c *Client) ExchangeCodeForToken(ctx context.Context, code string) (models.TokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s&code_verifier=%s",
		code,
		c.clientID,
		c.redirectURI,
		c.codeVerifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data))
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("executing token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Token endpoint error response: %s\n", string(body))
		return models.TokenResponse{}, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var token models.TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return models.TokenResponse{}, fmt.Errorf("decoding token response: %w", err)
	}

	// Validate token response
	if token.AccessToken == "" {
		return models.TokenResponse{}, fmt.Errorf("token response missing access_token")
	}

	if token.TokenType == "" {
		return models.TokenResponse{}, fmt.Errorf("token response missing token_type")
	}

	return token, nil
}

// RefreshAccessToken uses a refresh token to obtain a new access token
func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (models.TokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s&client_id=%s",
		refreshToken,
		c.clientID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data))
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("creating refresh token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("executing refresh token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.TokenResponse{}, fmt.Errorf("refresh token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var token models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return models.TokenResponse{}, fmt.Errorf("decoding refresh token response: %w", err)
	}

	return token, nil
}

func generateCodeVerifier() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
