package models

import "time"

// ParentResponse represents the full response from the parent API
type ParentResponse struct {
	TotalUsers   int       `json:"total_users"`
	TotalPosts   int       `json:"total_posts"`
	TotalFollows int       `json:"total_follows"`
	TotalLikes   int       `json:"total_likes"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CachedResponse represents our enhanced response with growth rate
type CachedResponse struct {
	TotalUsers               int       `json:"total_users"`
	TotalPosts               int       `json:"total_posts"`
	TotalFollows             int       `json:"total_follows"`
	TotalLikes               int       `json:"total_likes"`
	UsersGrowthRatePerSecond float64   `json:"users_growth_rate_per_second"`
	LastUpdateTime           time.Time `json:"last_update_time"`
	NextUpdateTime           time.Time `json:"next_update_time"`
}
