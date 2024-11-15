package api

import (
	"fmt"
	"net/http"

	"github.com/dickeyy/bsky-stats-api/config"
	"github.com/dickeyy/bsky-stats-api/internal/cache"
	"github.com/dickeyy/bsky-stats-api/internal/client"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg    *config.Config
	cache  *cache.Cache
	client *client.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		cache:  cache.NewCache(cfg.CacheExpiration),
		client: client.NewClient(cfg.ParentAPIURL),
	}
}

func (s *Server) Start() error {
	// Register handlers
	http.HandleFunc("/", s.handleStats)

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Info().
		Str("env", s.cfg.Env).
		Str("port", s.cfg.Port).
		Msg("Starting server")

	return http.ListenAndServe(addr, nil)
}
