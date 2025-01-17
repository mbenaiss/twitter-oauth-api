package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"myapp/handlers"
)

func main() {
	// Verify required environment variables
	if os.Getenv("CLIENT_ID") == "" || os.Getenv("CLIENT_SECRET") == "" {
		log.Fatal("CLIENT_ID and CLIENT_SECRET environment variables must be set")
	}

	// Initialize Gin
	router := gin.Default()

	// Setup session middleware with secure cookie store
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	router.Use(sessions.Sessions("oauth-session", store))

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Setup routes
	router.GET("/", handlers.HomeHandler)
	router.GET("/login", handlers.LoginHandler)
	router.GET("/callback", handlers.CallbackHandler)
	router.GET("/logout", handlers.LogoutHandler)

	// Start server on 0.0.0.0:8000
	log.Println("Server starting on http://0.0.0.0:8000")
	if err := router.Run("0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}
}