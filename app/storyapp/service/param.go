package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type AddStoryRequest struct {
	ID              types.ID  `json:"-"`
	CreatorID       types.ID  `json:"-"`
	MediaFileID     types.ID  `json:"media_file_id"`
	Title           string    `json:"title"`
	Caption         string    `json:"caption"`
	LinkURL         string    `json:"link_url"`
	LinkText        string    `json:"link_text"`
	DurationSeconds int       `json:"duration_seconds"`
	PublishAt       time.Time `json:"publish_at"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type AddStoryResponse struct {
	Story
}
