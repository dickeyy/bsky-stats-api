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
	Data      models.ParentResponse
	ExpiresAt time.Time
}

func NewCache(expiration time.Duration) *Cache {
	return &Cache{
		expiration: expiration,
	}
}

func (c *Cache) Get() (*models.ParentResponse, bool) {
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

func (c *Cache) Set(data models.ParentResponse) {
	c.Lock()
	defer c.Unlock()

	expiresAt := time.Now().Add(c.expiration)
	c.data = &CacheEntry{
		Data:      data,
		ExpiresAt: expiresAt,
	}

	log.Info().
		Time("expires_at", expiresAt).
		Int("total_users", data.TotalUsers).
		Int("total_posts", data.TotalPosts).
		Int("total_follows", data.TotalFollows).
		Int("total_likes", data.TotalLikes).
		Time("updated_at", data.UpdatedAt).
		Msg("Cache updated")
}
