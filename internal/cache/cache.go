package cache

import (
	"sync"
	"time"

	"github.com/dickeyy/bsky-stats-api/internal/models"

	"github.com/rs/zerolog/log"
)

type Cache struct {
	sync.RWMutex
	data       *CacheEntry
	expiration time.Duration
}

type CacheEntry struct {
	Data      models.CachedResponse
	ExpiresAt time.Time
}

func NewCache(expiration time.Duration) *Cache {
	return &Cache{
		expiration: expiration,
	}
}

func (c *Cache) Get() (*models.CachedResponse, bool) {
	c.RLock()
	defer c.RUnlock()

	if c.data == nil {
		log.Debug().Msg("Cache is empty")
		return nil, false
	}

	if time.Now().After(c.data.ExpiresAt) {
		log.Debug().
			Time("expires_at", c.data.ExpiresAt).
			Msg("Cache expired")
		return nil, false
	}

	timeLeft := time.Until(c.data.ExpiresAt)
	log.Debug().
		Dur("time_left", timeLeft).
		Time("expires_at", c.data.ExpiresAt).
		Msg("Cache hit")

	return &c.data.Data, true
}

func (c *Cache) GetPrevious() *models.CachedResponse {
	c.RLock()
	defer c.RUnlock()

	if c.data == nil {
		return nil
	}
	return &c.data.Data
}

func (c *Cache) Set(data models.CachedResponse) {
	c.Lock()
	defer c.Unlock()

	c.data = &CacheEntry{
		Data:      data,
		ExpiresAt: data.NextUpdateTime,
	}

	log.Info().
		Time("next_update", data.NextUpdateTime).
		Int("total_users", data.TotalUsers).
		Int("total_posts", data.TotalPosts).
		Int("total_follows", data.TotalFollows).
		Int("total_likes", data.TotalLikes).
		Float64("growth_rate", data.UsersGrowthRatePerSecond).
		Msg("Cache updated")
}
