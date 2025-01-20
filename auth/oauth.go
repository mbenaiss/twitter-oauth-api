package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	AuthURL      = "https://x.com/i/oauth2/authorize"
	TokenURL     = "https://api.x.com/2/oauth2/token"
	UserInfoURL  = "https://api.x.com/2/users/me"
	RedirectURI  = "http://localhost:8000/callback"
	OAuthScopes  = "tweet.read users.read follows.read offline.access"
)

type Client struct {
	clientID     string
	clientSecret string
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GenerateCodeVerifier creates a random code verifier for PKCE
func GenerateCodeVerifier() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateCodeChallenge creates SHA256 challenge from verifier
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// GenerateState creates a random state parameter for CSRF protection
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ExchangeCodeForToken exchanges auth code for access token
func (c *Client) ExchangeCodeForToken(code, verifier string) (*TokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s&code_verifier=%s",
		code,
		c.clientID,
		RedirectURI,
		verifier,
	)

	req, err := http.NewRequest(http.MethodPost, TokenURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return &token, nil
}

// RefreshAccessToken uses a refresh token to obtain a new access token
func (c *Client) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s&client_id=%s",
		refreshToken,
		c.clientID,
	)

	req, err := http.NewRequest(http.MethodPost, TokenURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating refresh token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing refresh token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding refresh token response: %w", err)
	}

	return &token, nil
}

// GetClientID returns the client ID
func (c *Client) GetClientID() string {
	return c.clientID
}