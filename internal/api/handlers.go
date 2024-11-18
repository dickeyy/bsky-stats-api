package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) handleStats(c *gin.Context) {
	logger := log.With().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Logger()

	// Try to get data from cache
	if data, ok := s.cache.Get(); ok {
		// Cache hit
		c.Header("X-Cache", "HIT")
		c.JSON(http.StatusOK, data)
		return
	}

	// Get previous cache for growth rate calculation
	prevCache := s.cache.GetPrevious()

	// Cache miss or expired, fetch and process new stats
	data, err := s.client.ProcessStats(prevCache)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to fetch data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update cache
	s.cache.Set(*data)

	// Return response
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, data)
}

func (s *Server) handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
