package service

import (
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
	Data              string                   `json:"data"`
	ChannelDeliveries []ChannelDeliveryRequest `json:"channel_deliveries"`
}

type SendNotificationResponse struct {
	Notification
}
