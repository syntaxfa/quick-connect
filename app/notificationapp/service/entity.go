package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

// Notification represents a single notification sent to a user.
// It tracks the content, recipient, and detailed delivery status across various channels.
type Notification struct {
	ID                types.ID          `json:"id"`
	UserID            types.ID          `json:"user_id"`
	Type              NotificationType  `json:"type"`
	Title             string            `json:"title"`
	Body              string            `json:"body"`
	Data              json.RawMessage   `json:"data,omitempty"`
	IsRead            bool              `json:"is_read"`
	CreatedAt         time.Time         `json:"created_at"`
	ChannelDeliveries []ChannelDelivery `json:"channel_deliveries"`
	OverallStatus     OverallStatus     `json:"overall_status"`
}

// OverallStatus defines the aggregate delivery status of a notification across all channels.
// This status reflects whether delivery to all intended channels was successful, failed, or is pending.
type OverallStatus uint8

const (
	OverallStatusPending  OverallStatus = iota + 1 // Notification is newly created and processing/delivery attempts are pending
	OverallStatusSent                              // All requested channels have successfully delivered the notification
	OverallStatusFailed                            // Delivery to all critical/requested channels failed after all retries
	OverallStatusRetrying                          // At least one channel is still in a retrying state
	OverallStatusIgnored                           // Some channels succeeded, while others failed or are still pending (partial success)
)

// NotificationType defines the categorization of the notification.
// This helps in distinguishing different kinds of messages and can be used for
// user preferences or specific business logic.
type NotificationType uint8

func NotificationTypeToInt(notificationType string) (NotificationType, error) {
	switch notificationType {
	case "optional":
		return NotificationTypeOptional, nil
	case "info":
		return NotificationTypeInfo, nil
	case "promotion":
		return NotificationTypePromotion, nil
	case "critical":
		return NotificationTypeCritical, nil
	default:
		return 0, errors.New("invalid notification type")
	}
}

func NotificationTypeToString(notificationType NotificationType) (string, error) {
	switch notificationType {
	case NotificationTypeOptional:
		return "optional", nil
	case NotificationTypeInfo:
		return "info", nil
	case NotificationTypePromotion:
		return "promotion", nil
	case NotificationTypeCritical:
		return "critical", nil
	default:
		return "", errors.New("invalid notification type")
	}
}

const (
	NotificationTypeOptional  NotificationType = iota + 1 // Can be opted out by user preferences
	NotificationTypeInfo                                  // General informational messages
	NotificationTypePromotion                             // Marketing or promotional messages
	NotificationTypeCritical                              // High-priority messages that usually cannot be opted out of (e.g., security alerts, password resets)
)

// ChannelType defines the various communication channels through which a notification can be sent.
type ChannelType uint8

const (
	ChannelTypeSMS     ChannelType = iota + 1 // Short Message Service (text messages)
	ChannelTypeEmail                          // Electronic mail
	ChannelTypeWebPush                        // Browser-based push notifications (e.g., via FCM, Web Push API)
)

// DeliveryStatus represents the status of a single delivery attempt for a specific channel.
type DeliveryStatus uint8

const (
	DeliveryStatusPending  DeliveryStatus = iota + 1 // Delivery for this channel has not been attempted yet, or is awaiting processing
	DeliveryStatusSent                               // The notification was successfully sent through this channel
	DeliveryStatusFailed                             // Delivery through this channel failed (after exhausting all retries)
	DeliveryStatusRetrying                           // Delivery for this channel is currently being retried
	DeliveryStatusIgnored                            // Delivery to this channel was ignored (e.g., user opted out, channel not configured, or a business rule prevented sending)
)

// ChannelDelivery represents the detailed status of a notification's delivery attempt for a specific channel.
type ChannelDelivery struct {
	Channel       ChannelType    `json:"channel"`         // The type of communication channel
	Status        DeliveryStatus `json:"status"`          // The current delivery status for this channel
	LastAttemptAt *time.Time     `json:"last_attempt_at"` // Timestamp of the last delivery attempt
	AttemptCount  int            `json:"attempt_count"`   // Number of times delivery has been attempted for this channel
	Error         *string        `json:"error"`           // Optional: Error message if the last delivery attempt failed
}
