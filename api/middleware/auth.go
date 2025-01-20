package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the request has a valid API key
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("api_key")
		if apiKey != secret || apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
