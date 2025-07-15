package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

// Notification represents a single notification sent to a user.
// It tracks the content, recipient, and detailed delivery status across various channels.
type Notification struct {
	ID                types.ID          `json:"id"`
	UserID            types.ID          `json:"user_id"`
	Type              NotificationType  `json:"type"`
	Data              map[string]string `json:"data,omitempty"`
	TemplateName      string            `json:"template_name"`
	DynamicBodyData   map[string]string `json:"dynamic_body_data,omitempty"`
	DynamicTitleData  map[string]string `json:"dynamic_title_data,omitempty"`
	IsRead            bool              `json:"is_read"`
	IsInApp           bool              `json:"is_in_app"`
	CreatedAt         time.Time         `json:"created_at"`
	ChannelDeliveries []ChannelDelivery `json:"channel_deliveries"`
	OverallStatus     OverallStatus     `json:"overall_status"`
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

func IsValidOverallStatus(overallStatus OverallStatus) bool {
	if overallStatus == OverallStatusPending || overallStatus == OverallStatusSent ||
		overallStatus == OverallStatusFailed || overallStatus == OverallStatusRetrying ||
		overallStatus == OverallStatusIgnored || overallStatus == OverallStatusMixed {
		return true
	}

	return false
}

// NotificationType defines the categorization of the notification.
// This helps in distinguishing different kinds of messages and can be used for
// user preferences or specific business logic.
type NotificationType string

const (
	NotificationTypeOptional  NotificationType = "optional"  // Can be opted out by user preferences
	NotificationTypeInfo      NotificationType = "info"      // General informational messages
	NotificationTypePromotion NotificationType = "promotion" // Marketing or promotional messages
	NotificationTypeCritical  NotificationType = "critical"  // High-priority messages that usually cannot be opted out of (e.g., security alerts, password resets)
	NotificationTypeDirect    NotificationType = "direct"    // Notifications of "direct" type do not have a retry mechanism in case of delivery failure.
)

func IsValidNotificationType(notificationType NotificationType) bool {
	if notificationType == NotificationTypeOptional || notificationType == NotificationTypeInfo ||
		notificationType == NotificationTypePromotion || notificationType == NotificationTypeCritical ||
		notificationType == NotificationTypeDirect {
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
	ChannelTypeInApp   ChannelType = "in_app"   // In-App notification
)

func IsValidChannelType(channelType ChannelType) bool {
	if channelType == ChannelTypeSMS || channelType == ChannelTypeEmail || channelType == ChannelTypeWebPush || channelType == ChannelTypeInApp {
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

// Template represents a notification template definition.
// It groups different content variations (bodies) for various channels and languages
// under a single logical template name.
type Template struct {
	ID        types.ID          `json:"id"`
	Name      string            `json:"name"`
	Contents  []TemplateContent `json:"contents"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// TemplateContent defines the content of a specific template for a given channel.
type TemplateContent struct {
	Channel ChannelType   `json:"channel"`
	Bodies  []ContentBody `json:"bodies"`
}

// ContentBody defines the content body of a specific template content for a given language.
type ContentBody struct {
	Lang  string `json:"lang"`
	Body  string `json:"body"`
	Title string `json:"title"`
}

// UserSetting A user can have their custom and personalized settings, such as language and channels
// they do not want to receive notifications from.
type UserSetting struct {
	ID             types.ID        `json:"id"`
	UserID         types.ID        `json:"user_id"`
	Lang           string          `json:"lang"`
	IgnoreChannels []IgnoreChannel `json:"ignore_channels"`
}

// IgnoreChannel A user can ignore channels with a high level of customization. A user can specify based on notification type,
// for example, only promotion notifications should be ignored in SMS.
// Note: A user cannot ignore notifications whose type is critical or direct.
type IgnoreChannel struct {
	Channel           ChannelType        `json:"channel"`
	NotificationTypes []NotificationType `json:"notification_type"`
}
