package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	logger := log.With().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Logger()

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Handle OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		logger.Warn().Msg("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	// Try to get data from cache
	if data, ok := s.cache.Get(); ok {
		// Cache hit
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(data)

		logger.Info().
			Str("cache", "HIT").
			Dur("duration", time.Since(start)).
			Msg("Request completed")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update cache
	s.cache.Set(*data)

	// Return response
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(data)

	logger.Info().
		Str("cache", "MISS").
		Dur("duration", time.Since(start)).
		Float64("growth_rate", data.UsersGrowthRatePerSecond).
		Msg("Request completed")
}
