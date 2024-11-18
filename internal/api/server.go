package api

import (
	"fmt"
	"net/http"

	"github.com/dickeyy/bsky-stats-api/config"
	"github.com/dickeyy/bsky-stats-api/internal/cache"
	"github.com/dickeyy/bsky-stats-api/internal/client"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	cache  *cache.Cache
	client *client.Client
}

func NewServer(cfg *config.Config) *Server {
	router := gin.New()

	// Use gin.Recovery() middleware to handle panics
	router.Use(gin.Recovery())

	// Add zerolog logger middleware
	router.Use(loggerMiddleware())

	// Add CORS middleware
	router.Use(corsMiddleware())

	server := &Server{
		router: router,
		cfg:    cfg,
		cache:  cache.NewCache(cfg.CacheExpiration),
		client: client.NewClient(cfg.ParentAPIURL),
	}

	// Register routes
	server.registerRoutes()

	return server
}

func (s *Server) registerRoutes() {
	s.router.GET("/", s.handleStats)
	s.router.GET("/ping", s.handlePing)
	s.router.HEAD("/ping", s.handlePing)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Info().
		Str("env", s.cfg.Env).
		Str("port", s.cfg.Port).
		Msg("Starting server")

	return s.router.Run(addr)
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log the request completion
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Msg("Request processed")
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
