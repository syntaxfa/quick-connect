package http

import (
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
		ID:    claims.ID,
		Roles: roles,
	}
}
