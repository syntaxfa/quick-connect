package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type Story struct {
	ID              types.ID  `json:"id"`
	MediaFileID     types.ID  `json:"media_file_id"`
	Title           string    `json:"title"`
	Caption         string    `json:"caption"`
	LinkURL         string    `json:"link_url"`
	LinkText        string    `json:"link_text"`
	DurationSeconds int       `json:"duration_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	PublishAt       time.Time `json:"publish_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatorID       types.ID  `json:"creator_id"`
	IsViewed        bool      `json:"is_viewed"`
}

func (s *Story) IsVisible() bool {
	now := time.Now()

	return s.IsActive &&
		(now.After(s.PublishAt) || now.Equal(s.PublishAt)) &&
		now.Before(s.ExpiresAt)
}

type StoryView struct {
	StoryID  types.ID
	ViewerID types.ID
	ViewedAt time.Time
}
