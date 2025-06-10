package service

import (
	"encoding/json"

	"github.com/syntaxfa/quick-connect/types"
)

type ChannelDeliveryRequest struct {
	Channel ChannelType `json:"channel"`
}

type SendNotificationRequest struct {
	UserID            types.ID                 `json:"user_id"`
	ExternalUserID    string                   `json:"-"`
	Type              NotificationType         `json:"type"`
	Title             string                   `json:"title"`
	Body              string                   `json:"body"`
	Data              json.RawMessage          `json:"data"`
	ChannelDeliveries []ChannelDeliveryRequest `json:"channel_deliveries"`
}

type SendNotificationResponse struct {
	Notification
}
