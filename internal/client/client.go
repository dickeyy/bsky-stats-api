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

func (c *Client) FetchStats() (*models.ParentResponse, error) {
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
