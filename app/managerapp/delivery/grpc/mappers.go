package grpc

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
)

func convertUserRoleToPB(roles []types.Role) []userpb.Role {
	var pbRoles []userpb.Role
	for _, role := range roles {
		switch role {
		case types.RoleSuperUser:
			pbRoles = append(pbRoles, userpb.Role_ROLE_SUPERUSER)
		case types.RoleSupport:
			pbRoles = append(pbRoles, userpb.Role_ROLE_SUPPORT)
		case types.RoleStory:
			pbRoles = append(pbRoles, userpb.Role_ROLE_STORY)
		case types.RoleFile:
			pbRoles = append(pbRoles, userpb.Role_ROLE_FILE)
		case types.RoleNotification:
			pbRoles = append(pbRoles, userpb.Role_ROLE_NOTIFICATION)
		}
	}

	return pbRoles
}

func convertUserRoleToEntity(pbRoles []userpb.Role) []types.Role {
	var roles []types.Role
	for _, role := range pbRoles {
		switch role {
		case userpb.Role_ROLE_SUPERUSER:
			roles = append(roles, types.RoleSuperUser)
		case userpb.Role_ROLE_SUPPORT:
			roles = append(roles, types.RoleSupport)
		case userpb.Role_ROLE_STORY:
			roles = append(roles, types.RoleStory)
		case userpb.Role_ROLE_FILE:
			roles = append(roles, types.RoleFile)
		case userpb.Role_ROLE_NOTIFICATION:
			roles = append(roles, types.RoleNotification)
		}
	}

	return roles
}

func convertCreateUserRequestToEntity(req *userpb.CreateUserRequest) userservice.UserCreateRequest {
	return userservice.UserCreateRequest{
		ID:          "",
		Username:    req.Username,
		Password:    req.Password,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Roles:       convertUserRoleToEntity(req.Roles),
	}
}

func convertUserToPB(req userservice.User) *userpb.User {
	return &userpb.User{
		Id:           string(req.ID),
		Username:     req.Username,
		Fullname:     req.Fullname,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Avatar:       req.Avatar,
		Roles:        convertUserRoleToPB(req.Roles),
		LastOnlineAt: nil,
	}
}

func convertUserUpdateFromSuperuserToEntity(req *userpb.UserUpdateFromSuperUserRequest) userservice.UserUpdateFromSuperuserRequest {
	return userservice.UserUpdateFromSuperuserRequest{
		Username:    req.Username,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Roles:       convertUserRoleToEntity(req.Roles),
	}
}

func convertUserUpdateFromOwnToEntity(req *userpb.UserUpdateFromOwnRequest) userservice.UserUpdateFromOwnRequest {
	return userservice.UserUpdateFromOwnRequest{
		Username:    req.Username,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
}

func convertUserListRequestToEntity(req *userpb.UserListRequest) userservice.ListUserRequest {
	request := userservice.ListUserRequest{
		Username: req.Username,
		Paginated: paginate.RequestBase{
			CurrentPage: req.CurrentPage,
			PageSize:    req.PageSize,
		},
	}

	switch req.SortDirection {
	case userpb.SortDirection_SORT_DIRECTION_ASC:
		request.Paginated.Descending = false
	default:
		request.Paginated.Descending = true
	}

	return request
}

func convertUserListResponseToPB(req userservice.ListUserResponse) *userpb.UserListResponse {
	var users []*userpb.User
	for _, user := range req.Results {
		users = append(users, convertUserToPB(user))
	}

	return &userpb.UserListResponse{
		CurrentPage: req.Paginate.CurrentPage,
		PageSize:    req.Paginate.PageSize,
		TotalNumber: req.Paginate.TotalNumbers,
		TotalPage:   req.Paginate.TotalPage,
		Users:       users,
	}
}
