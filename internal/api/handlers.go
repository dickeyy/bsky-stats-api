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
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(data)

		logger.Info().
			Str("cache", "HIT").
			Dur("duration", time.Since(start)).
			Msg("Request completed")
		return
	}

	// Cache miss or expired, fetch from parent API
	data, err := s.client.FetchStats()
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(data)

	logger.Info().
		Str("cache", "MISS").
		Dur("duration", time.Since(start)).
		Msg("Request completed")
}
