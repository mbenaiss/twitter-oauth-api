package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mbenaiss/twitter-oauth-api/api/middleware"
	"github.com/mbenaiss/twitter-oauth-api/twitter"
)

// Server is the main server struct
type Server struct {
	router *gin.Engine
	port   string
	secret string
}

// NewServer creates a new server instance
func NewServer(port string, secret string) *Server {
	router := gin.Default()

	return &Server{
		router: router,
		port:   port,
		secret: secret,
	}
}

// Start starts the server
func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes(authClient *twitter.Client) {
	authMiddleware := middleware.AuthMiddleware(s.secret)
	s.router.Use(authMiddleware).GET("/", loginHandler(authClient))
	s.router.Use(authMiddleware).POST("/refresh", refreshTokenHandler(authClient))

	s.router.GET("/callback", callbackHandler(authClient))
}
