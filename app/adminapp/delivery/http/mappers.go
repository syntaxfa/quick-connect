package http

import (
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
)

type User struct {
	ID           string
	Username     string
	Fullname     string
	Email        string
	PhoneNumber  string
	Avatar       string
	Roles        []string
	LastOnlineAt string
}

// RoleInfo struct helper for templates.
type RoleInfo struct {
	Name  string // "SUPERUSER", "SUPPORT"
	Value int32  // 1, 2
}

// GetAllRoles returns all available roles for the template.
func GetAllRoles() []RoleInfo {
	return []RoleInfo{
		{Name: string(types.RoleSuperUser), Value: int32(userpb.Role_ROLE_SUPERUSER)},
		{Name: string(types.RoleSupport), Value: int32(userpb.Role_ROLE_SUPPORT)},
		{Name: string(types.RoleStory), Value: int32(userpb.Role_ROLE_STORY)},
		{Name: string(types.RoleFile), Value: int32(userpb.Role_ROLE_FILE)},
		{Name: string(types.RoleNotification), Value: int32(userpb.Role_ROLE_NOTIFICATION)},
		{Name: string(types.RoleClient), Value: int32(userpb.Role_ROLE_CLIENT)},
		{Name: string(types.RoleGuest), Value: int32(userpb.Role_ROLE_GUEST)},
		{Name: string(types.RoleBot), Value: int32(userpb.Role_ROLE_BOT)},
	}
}

// HasRole checks if a user has a specific role (used by template).
func HasRole(userRoles []string, targetRole string) bool {
	for _, r := range userRoles {
		if r == targetRole {
			return true
		}
	}
	return false
}

// ParseRolesFromForm converts string roles from form back to gRPC Enum.
func ParseRolesFromForm(roleStrings []string) []userpb.Role {
	var roles []userpb.Role
	roleMap := map[string]userpb.Role{
		string(types.RoleSuperUser):    userpb.Role_ROLE_SUPERUSER,
		string(types.RoleSupport):      userpb.Role_ROLE_SUPPORT,
		string(types.RoleStory):        userpb.Role_ROLE_STORY,
		string(types.RoleFile):         userpb.Role_ROLE_FILE,
		string(types.RoleNotification): userpb.Role_ROLE_NOTIFICATION,
		string(types.RoleClient):       userpb.Role_ROLE_CLIENT,
		string(types.RoleGuest):        userpb.Role_ROLE_GUEST,
		string(types.RoleBot):          userpb.Role_ROLE_BOT,
	}

	for _, rs := range roleStrings {
		if val, ok := roleMap[rs]; ok {
			roles = append(roles, val)
		}
	}
	return roles
}

func convertUserPbToUser(userPb *userpb.User) User {
	var roles []string
	for _, role := range userPb.GetRoles() {
		switch role {
		case userpb.Role_ROLE_SUPERUSER:
			roles = append(roles, string(types.RoleSuperUser))
		case userpb.Role_ROLE_SUPPORT:
			roles = append(roles, string(types.RoleSupport))
		case userpb.Role_ROLE_NOTIFICATION:
			roles = append(roles, string(types.RoleNotification))
		case userpb.Role_ROLE_STORY:
			roles = append(roles, string(types.RoleStory))
		case userpb.Role_ROLE_FILE:
			roles = append(roles, string(types.RoleFile))
		case userpb.Role_ROLE_CLIENT:
			roles = append(roles, string(types.RoleClient))
		case userpb.Role_ROLE_GUEST:
			roles = append(roles, string(types.RoleGuest))
		case userpb.Role_ROLE_BOT:
			roles = append(roles, string(types.RoleBot))
		case userpb.Role_ROLE_SERVICE:
			roles = append(roles, string(types.RoleService))
		case userpb.Role_ROLE_UNSPECIFIED:
			continue
		}
	}

	lastOnline := ""
	if userPb.GetLastOnlineAt() != nil {
		t := userPb.GetLastOnlineAt().AsTime()
		lastOnline = t.Format("Jan 02, 2006 at 3:04 PM")
	}

	return User{
		ID:           userPb.GetId(),
		Username:     userPb.GetUsername(),
		Fullname:     userPb.GetFullname(),
		Email:        userPb.GetEmail(),
		PhoneNumber:  userPb.GetPhoneNumber(),
		Avatar:       userPb.GetAvatar(),
		Roles:        roles,
		LastOnlineAt: lastOnline,
	}
}

func convertClaimsToUser(claims *types.UserClaims) User {
	var roles []string
	for _, role := range claims.Roles {
		roles = append(roles, string(role))
	}

	return User{
		ID:    string(claims.UserID),
		Roles: roles,
	}
}

type Conversation struct {
	ID                  string
	ClientUserID        string
	AssignedSupportID   string
	Status              string
	LastMessageSnippet  string
	LastMessageSenderID string
	UpdatedAt           string
}

// ConversationStatusInfo struct helper for templates.
type ConversationStatusInfo struct {
	Name  string // "OPEN", "CLOSED"
	Value string // "open", "closed" (value for query param)
	PbVal int32  // Protobuf enum value
}

// GetAllConversationStatuses returns filterable statuses for the template.
func GetAllConversationStatuses() []ConversationStatusInfo {
	return []ConversationStatusInfo{
		{Name: "Open", Value: string(service.ConversationStatusOpen), PbVal: int32(conversationpb.Status_STATUS_OPEN)},
		{Name: "Closed", Value: string(service.ConversationStatusClosed), PbVal: int32(conversationpb.Status_STATUS_CLOSED)},
	}
}

// ParseStatusesFromForm converts string statuses from form back to gRPC Enum.
func ParseStatusesFromForm(statusStrings []string) []conversationpb.Status {
	var statuses []conversationpb.Status
	statusMap := map[string]conversationpb.Status{
		string(service.ConversationStatusOpen):        conversationpb.Status_STATUS_OPEN,
		string(service.ConversationStatusClosed):      conversationpb.Status_STATUS_CLOSED,
		string(service.ConversationStatusBotHandling): conversationpb.Status_STATUS_BOT_HANDLING,
		string(service.ConversationStatusNew):         conversationpb.Status_STATUS_NEW,
	}

	for _, ss := range statusStrings {
		if val, ok := statusMap[ss]; ok {
			statuses = append(statuses, val)
		}
	}
	return statuses
}

// convertConversationPbToConversation maps the gRPC Conversation object to the template-friendly Conversation struct.
func convertConversationPbToConversation(convPb *conversationpb.Conversation) Conversation {
	updatedAt := ""
	if convPb.GetUpdatedAt() != nil {
		t := convPb.GetUpdatedAt().AsTime()
		updatedAt = t.Format("3:04 PM")
	}

	return Conversation{
		ID:                  convPb.GetId(),
		ClientUserID:        convPb.GetClientUserId(),
		AssignedSupportID:   convPb.GetAssignedSupportId(),
		Status:              convertStatusPbToString(convPb.GetStatus()),
		LastMessageSnippet:  convPb.GetLastMessageSnippet(),
		LastMessageSenderID: convPb.GetLastMessageSenderId(),
		UpdatedAt:           updatedAt,
	}
}

func convertStatusPbToString(statusPb conversationpb.Status) string {
	switch statusPb {
	case conversationpb.Status_STATUS_NEW:
		return string(service.ConversationStatusNew)
	case conversationpb.Status_STATUS_OPEN:
		return string(service.ConversationStatusOpen)
	case conversationpb.Status_STATUS_CLOSED:
		return string(service.ConversationStatusClosed)
	case conversationpb.Status_STATUS_BOT_HANDLING:
		return string(service.ConversationStatusBotHandling)
	case conversationpb.Status_STATUS_UNSPECIFIED:
		return "unknown"
	default:
		return "unknown"
	}
}
