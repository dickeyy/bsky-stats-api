package main

import (
	"github.com/dickeyy/bsky-stats-api/config"
	"github.com/dickeyy/bsky-stats-api/internal/api"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize configuration
	cfg := config.Init()

	// Create and start server
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
