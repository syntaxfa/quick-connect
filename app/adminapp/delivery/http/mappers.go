package http

import (
	"time"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
)

func convertUserPbToUser(userPb *userpb.User) User {
	var roles []string
	for _, role := range userPb.Roles {
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
		}
	}

	return User{
		ID:           userPb.Id,
		Username:     userPb.Username,
		Fullname:     userPb.Fullname,
		Email:        userPb.Email,
		PhoneNumber:  userPb.PhoneNumber,
		Avatar:       userPb.Avatar,
		Roles:        roles,
		LastOnlineAt: time.Now().Add(-15 * time.Minute).Format("Jan 02, 2006 at 3:04 PM"),
	}
}
