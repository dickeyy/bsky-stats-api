package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dickeyy/bsky-stats-api/internal/models"
	"github.com/rs/zerolog/log"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) ProcessStats(prevCache *models.CachedResponse) (*models.CachedResponse, error) {
	parentResp, err := c.fetchFromParentAPI()
	if err != nil {
		return nil, err
	}

	// Calculate growth rate
	var usersGrowthRatePerSecond float64 = 2.34 * 60 * 0.95 // Default rate if no previous cache

	if prevCache != nil {
		timeSinceLast := parentResp.UpdatedAt.Sub(prevCache.LastUpdateTime).Seconds()
		if timeSinceLast > 0 {
			usersChangeSinceLast := float64(parentResp.TotalUsers - prevCache.TotalUsers)
			usersGrowthRatePerSecond = (usersChangeSinceLast / timeSinceLast) * 0.95
		}
	}

	// Create enhanced response
	response := &models.CachedResponse{
		TotalUsers:               parentResp.TotalUsers,
		TotalPosts:               parentResp.TotalPosts,
		TotalFollows:             parentResp.TotalFollows,
		TotalLikes:               parentResp.TotalLikes,
		UsersGrowthRatePerSecond: usersGrowthRatePerSecond,
		LastUpdateTime:           parentResp.UpdatedAt,
		NextUpdateTime:           time.Now().Add(time.Minute),
	}

	return response, nil
}

func (c *Client) fetchFromParentAPI() (*models.ParentResponse, error) {
	log.Debug().Msg("Fetching from parent API")

	start := time.Now()
	resp, err := c.httpClient.Get(c.baseURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", c.baseURL).
			Msg("Failed to fetch from parent API")
		return nil, fmt.Errorf("failed to fetch from parent API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var parentResp models.ParentResponse
	if err := json.Unmarshal(body, &parentResp); err != nil {
		log.Error().
			Err(err).
			Str("body", string(body)).
			Msg("Failed to unmarshal response")
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	log.Debug().
		Dur("duration", time.Since(start)).
		Int("status", resp.StatusCode).
		Int("body_size", len(body)).
		Msg("Parent API request completed")

	return &parentResp, nil
}
