package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"myapp/auth"
	"myapp/models"
)

// HomeHandler renders the home page
func HomeHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.HTML(http.StatusOK, "home.html", gin.H{
		"user": user,
	})
}

// LoginHandler initiates the OAuth2 PKCE flow
func LoginHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Generate and store PKCE verifier
	verifier, err := auth.GenerateCodeVerifier()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating verifier")
		return
	}
	challenge := auth.GenerateCodeChallenge(verifier)

	// Generate and store state
	state, err := auth.GenerateState()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating state")
		return
	}

	// Store in session
	session.Set("code_verifier", verifier)
	session.Set("state", state)
	session.Save()

	// Build authorization URL
	authURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
		auth.AuthURL,
		os.Getenv("CLIENT_ID"),
		auth.RedirectURI,
		auth.OAuthScopes,
		state,
		challenge,
	)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// CallbackHandler handles the OAuth2 callback
func CallbackHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Verify state parameter
	if state := c.Query("state"); state != session.Get("state") {
		c.String(http.StatusBadRequest, "Invalid state parameter")
		return
	}

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing authorization code")
		return
	}

	// Get stored verifier
	verifier := session.Get("code_verifier")
	if verifier == nil {
		c.String(http.StatusBadRequest, "Missing code verifier")
		return
	}

	// Exchange code for token
	token, err := auth.ExchangeCodeForToken(code, verifier.(string))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error exchanging code for token")
		return
	}

	// Get user info
	user, err := models.GetUserInfo(token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting user info")
		return
	}

	// Store user in session
	session.Set("user", user)
	session.Set("access_token", token.AccessToken)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// LogoutHandler clears the session
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}