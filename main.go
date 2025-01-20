package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"myapp/handlers"
	"myapp/auth"
)

type Config struct {
	TwitterClientID     string
	TwitterClientSecret string
	SessionSecret      string
}

func main() {
	// Initialize configuration
	config := Config{
		TwitterClientID:     os.Getenv("TWITTER_CLIENT_ID"),
		TwitterClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
	}

	// Verify required environment variables
	if config.TwitterClientID == "" || config.TwitterClientSecret == "" {
		log.Fatal("TWITTER_CLIENT_ID and TWITTER_CLIENT_SECRET environment variables must be set")
	}

	// Initialize auth client
	authClient := auth.NewClient(config.TwitterClientID, config.TwitterClientSecret)

	// Initialize Gin
	router := gin.Default()

	// Setup session middleware with secure cookie store
	store := cookie.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	router.Use(sessions.Sessions("oauth-session", store))

	// Add token refresh middleware
	router.Use(handlers.TokenMiddleware(authClient))

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Setup routes with auth client
	router.GET("/", handlers.HomeHandler)
	router.GET("/login", handlers.LoginHandler(authClient))
	router.GET("/callback", handlers.CallbackHandler(authClient))
	router.GET("/logout", handlers.LogoutHandler)
	router.POST("/refresh", handlers.RefreshTokenHandler(authClient))

	// Start server on 0.0.0.0:8000
	log.Println("Server starting on http://0.0.0.0:8000")
	if err := router.Run("0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}
}