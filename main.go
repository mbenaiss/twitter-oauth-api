package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbenaiss/twitter-oauth-api/api"
	"github.com/mbenaiss/twitter-oauth-api/config"
	"github.com/mbenaiss/twitter-oauth-api/twitter"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	authClient := twitter.NewClient(cfg.TwitterClientID, cfg.TwitterClientSecret, fmt.Sprintf("http://%s/callback", serverAddr))

	server := api.NewServer(cfg.Port, cfg.APIKey)
	server.SetupRoutes(authClient)
	
	log.Printf("Server starting on http://%s", serverAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
