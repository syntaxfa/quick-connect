package service

import (
	"encoding/json"
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
type OverallStatus string

const (
	OverallStatusPending  OverallStatus = "pending"  // Notification is newly created and processing/delivery attempts are pending
	OverallStatusSent     OverallStatus = "sent"     // All requested channels have successfully delivered the notification
	OverallStatusFailed   OverallStatus = "failed"   // Delivery to all critical/requested channels failed after all retries
	OverallStatusRetrying OverallStatus = "retrying" // At least one channel is still in a retrying state
	OverallStatusIgnored  OverallStatus = "ignored"  // Some channels succeeded, while others failed or are still pending (partial success)
	OverallStatusMixed    OverallStatus = "mixed"
)

func IsValidOverallStatus(overallStatus string) bool {
	if overallStatus == "pending" || overallStatus == "sent" || overallStatus == "failed" ||
		overallStatus == "retrying" || overallStatus == "ignored" {
		return true
	}

	return false
}

// NotificationType defines the categorization of the notification.
// This helps in distinguishing different kinds of messages and can be used for
// user preferences or specific business logic.
type NotificationType string

const (
	NotificationTypeOptional  NotificationType = "optional " // Can be opted out by user preferences
	NotificationTypeInfo      NotificationType = "info"      // General informational messages
	NotificationTypePromotion NotificationType = "promotion" // Marketing or promotional messages
	NotificationTypeCritical  NotificationType = "critical"  // High-priority messages that usually cannot be opted out of (e.g., security alerts, password resets)
)

func IsValidNotificationType(notificationType string) bool {
	if notificationType == "optional" || notificationType == "info" ||
		notificationType == "promotion" || notificationType == "critical" {
		return true
	}

	return false
}

// ChannelType defines the various communication channels through which a notification can be sent.
type ChannelType string

const (
	ChannelTypeSMS     ChannelType = "sms"      // Short Message Service (text messages)
	ChannelTypeEmail   ChannelType = "email"    // Electronic mail
	ChannelTypeWebPush ChannelType = "web_push" // Browser-based push notifications (e.g., via FCM, Web Push API)
)

func IsValidChannelType(channelType string) bool {
	if channelType == "sms" || channelType == "email" || channelType == "web_push" {
		return true
	}

	return false
}

// DeliveryStatus represents the status of a single delivery attempt for a specific channel.
type DeliveryStatus string

const (
	DeliveryStatusPending  DeliveryStatus = "pending"  // Delivery for this channel has not been attempted yet, or is awaiting processing
	DeliveryStatusSent     DeliveryStatus = "sent"     // The notification was successfully sent through this channel
	DeliveryStatusFailed   DeliveryStatus = "failed"   // Delivery through this channel failed (after exhausting all retries)
	DeliveryStatusRetrying DeliveryStatus = "retrying" // Delivery for this channel is currently being retried
	DeliveryStatusIgnored  DeliveryStatus = "ignored"  // Delivery to this channel was ignored (e.g., user opted out, channel not configured, or a business rule prevented sending)
)

// ChannelDelivery represents the detailed status of a notification's delivery attempt for a specific channel.
type ChannelDelivery struct {
	Channel       ChannelType    `json:"channel"`         // The type of communication channel
	Status        DeliveryStatus `json:"status"`          // The current delivery status for this channel
	LastAttemptAt *time.Time     `json:"last_attempt_at"` // Timestamp of the last delivery attempt
	AttemptCount  int            `json:"attempt_count"`   // Number of times delivery has been attempted for this channel
	Error         *string        `json:"error"`           // Optional: Error message if the last delivery attempt failed
}
