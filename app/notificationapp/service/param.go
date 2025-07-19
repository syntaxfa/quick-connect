package service

import (
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
	Data              map[string]string        `json:"data"`
	TemplateName      string                   `json:"template_name"`
	DynamicBodyData   map[string]string        `json:"dynamic_body_data,omitempty"`
	DynamicTitleData  map[string]string        `json:"dynamic_title_data,omitempty"`
	ChannelDeliveries []ChannelDeliveryRequest `json:"channel_deliveries"`
	IsInApp           bool                     `json:"-"`
	Status            DeliveryStatus           `json:"-"`
}

// NotificationMessage rendered notification.
type NotificationMessage struct {
	ID        types.ID          `json:"id"`
	UserID    types.ID          `json:"user_id"`
	Type      NotificationType  `json:"type"`
	Data      map[string]string `json:"data"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	IsRead    bool              `json:"is_read"`
	Timestamp int64             `json:"timestamp"`
}

type ListNotificationRequest struct {
	ExternalUserID string               `json:"-"`
	IsRead         *bool                `json:"is_read"`
	Paginated      paginate.RequestBase `json:"paginated"`
}

type ListNotificationResponse struct {
	Results  []NotificationMessage `json:"results"`
	Paginate paginate.ResponseBase `json:"paginate"`
}

type AddTemplateRequest struct {
	ID       types.ID          `json:"-"`
	Name     string            `json:"name"` // maximum is 255 characters.
	Contents []TemplateContent `json:"contents"`
}

type UpdateUserSettingRequest struct {
	Lang           string          `json:"lang"`
	IgnoreChannels []IgnoreChannel `json:"ignore_channels"`
}

type ListTemplateRequest struct {
	Name      string               `json:"template_name"`
	Paginated paginate.RequestBase `json:"paginated"`
}

type ListTemplateResult struct {
	ID        types.ID  `json:"id"`
	Name      string    `json:"template_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListTemplateResponse struct {
	Results  []ListTemplateResult  `json:"results"`
	Paginate paginate.ResponseBase `json:"paginate"`
}
