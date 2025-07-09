package service

import (
	"encoding/json"
	"time"

	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/types"
)

type ChannelDeliveryRequest struct {
	Channel ChannelType `json:"channel"`
}

type SendNotificationRequest struct {
	ID                types.ID                 `json:"-"`
	UserID            types.ID                 `json:"-"`
	ExternalUserID    string                   `json:"external_user_id"`
	Type              NotificationType         `json:"type"`
	Title             string                   `json:"title"`
	Body              string                   `json:"body"`
	Data              json.RawMessage          `json:"data"`
	ChannelDeliveries []ChannelDeliveryRequest `json:"channel_deliveries"`
}

type SendNotificationRequestSchema struct {
	ID                types.ID                 `json:"-"`
	UserID            types.ID                 `json:"-"`
	ExternalUserID    string                   `json:"external_user_id"`
	Type              NotificationType         `json:"type"`
	Title             string                   `json:"title"`
	Body              string                   `json:"body"`
	Data              string                   `json:"data"`
	ChannelDeliveries []ChannelDeliveryRequest `json:"channel_deliveries"`
}

type SendNotificationResponse struct {
	Notification
}

type SendNotificationResponseSchema struct {
	ID                types.ID          `json:"id"`
	UserID            types.ID          `json:"user_id"`
	Type              NotificationType  `json:"type"`
	Title             string            `json:"title"`
	Body              string            `json:"body"`
	Data              string            `json:"data,omitempty"`
	IsRead            bool              `json:"is_read"`
	CreatedAt         time.Time         `json:"created_at"`
	ChannelDeliveries []ChannelDelivery `json:"channel_deliveries"`
	OverallStatus     OverallStatus     `json:"overall_status"`
}

type ListNotificationRequest struct {
	ExternalUserID string               `json:"-"`
	IsRead         *bool                `json:"is_read"`
	Paginated      paginate.RequestBase `json:"paginated"`
}

type ListNotificationResult struct {
	ID        types.ID         `json:"id"`
	UserID    types.ID         `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	Data      json.RawMessage  `json:"data,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type ListNotificationResultSchema struct {
	ID        types.ID         `json:"id"`
	UserID    types.ID         `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	Data      string           `json:"data,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type ListNotificationResponse struct {
	Results  []ListNotificationResult `json:"results"`
	Paginate paginate.ResponseBase    `json:"paginate"`
}

type ListNotificationResponseSchema struct {
	Results  []ListNotificationResultSchema `json:"results"`
	Paginate paginate.ResponseBase          `json:"paginate"`
}

type AddTemplateRequest struct {
	Name   string         `json:"name"` // maximum is 255 characters.
	Bodies []TemplateBody `json:"bodies"`
}

type UpdateUserSettingRequest struct {
	Lang           string          `json:"lang"`
	IgnoreChannels []IgnoreChannel `json:"ignore_channels"`
}
