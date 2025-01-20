package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mbenaiss/twitter-oauth-api/twitter"
)

func loginHandler(authClient *twitter.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authURL, err := authClient.GetAuthURL()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error getting auth URL: " + err.Error(),
			})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

func callbackHandler(authClient *twitter.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for error parameter
		if errMsg := c.Query("error"); errMsg != "" {
			errDesc := c.Query("error_description")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Authorization failed: %s - %s", errMsg, errDesc),
			})
			return
		}

		// Get required parameters from query
		state := c.Query("state")
		if state == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing state parameter",
			})
			return
		}

		// Get authorization code
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing authorization code",
			})
			return
		}

		token, err := authClient.ExchangeCodeForToken(c.Request.Context(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error exchanging code for token",
			})
			return
		}

		// Get user info
		user, err := authClient.GetUserInfo(c.Request.Context(), token.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error getting user info: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token,
		})
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func refreshTokenHandler(authClient *twitter.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req refreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		refreshToken := req.RefreshToken
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
			return
		}

		token, err := authClient.RefreshAccessToken(context.Background(), refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
