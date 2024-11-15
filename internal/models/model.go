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
