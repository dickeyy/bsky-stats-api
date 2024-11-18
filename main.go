package main

import (
	"github.com/dickeyy/bsky-stats-api/config"
	"github.com/dickeyy/bsky-stats-api/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize configuration
	cfg := config.Init()

	// Set Gin mode based on environment
	if cfg.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create and start server
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
