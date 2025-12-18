package manager

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserInfoToPB(resp userservice.UserInfoResponse) *userinternalpb.UserInfoResponse {
	return &userinternalpb.UserInfoResponse{
		Id:           string(resp.ID),
		Fullname:     resp.Fullname,
		Username:     resp.Username,
		Email:        resp.Email,
		PhoneNumber:  resp.PhoneNumber,
		Avatar:       resp.Avatar,
		LastOnlineAt: timestamppb.New(resp.LastOnlineAt),
	}
}

func convertLoginRequestToEntity(req *authpb.LoginRequest) userservice.UserLoginRequest {
	return userservice.UserLoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
}

func convertLoginResponseToPB(resp userservice.UserLoginResponse) *authpb.LoginResponse {
	return &authpb.LoginResponse{
		AccessToken:      resp.Token.AccessToken,
		RefreshToken:     resp.Token.RefreshToken,
		AccessExpiresIn:  resp.Token.AccessExpiresIn,
		RefreshExpiresIn: resp.Token.RefreshExpireIn,
	}
}

func convertTokenGenerateResponseToPB(resp *tokenservice.TokenGenerateResponse) *authpb.TokenRefreshResponse {
	return &authpb.TokenRefreshResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpireIn,
	}
}

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
		case types.RoleClient:
			pbRoles = append(pbRoles, userpb.Role_ROLE_CLIENT)
		case types.RoleGuest:
			pbRoles = append(pbRoles, userpb.Role_ROLE_GUEST)
		case types.RoleBot:
			pbRoles = append(pbRoles, userpb.Role_ROLE_BOT)
		case types.RoleService:
			pbRoles = append(pbRoles, userpb.Role_ROLE_SERVICE)
		}
	}

	return pbRoles
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
		LastOnlineAt: timestamppb.New(req.LastOnlineAt),
	}
}

func convertUserUpdateFromOwnToEntity(req *userpb.UserUpdateFromOwnRequest) userservice.UserUpdateFromOwnRequest {
	return userservice.UserUpdateFromOwnRequest{
		Username:    req.GetUsername(),
		Fullname:    req.GetFullname(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
	}
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
		case userpb.Role_ROLE_CLIENT:
			roles = append(roles, types.RoleClient)
		case userpb.Role_ROLE_GUEST:
			roles = append(roles, types.RoleGuest)
		case userpb.Role_ROLE_BOT:
			roles = append(roles, types.RoleBot)
		case userpb.Role_ROLE_SERVICE:
			roles = append(roles, types.RoleService)
		case userpb.Role_ROLE_UNSPECIFIED:
			continue
		}
	}

	return roles
}

func convertUserListRequestToEntity(req *userpb.UserListRequest) userservice.ListUserRequest {
	request := userservice.ListUserRequest{
		Username: req.GetUsername(),
		Paginated: paginate.RequestBase{
			CurrentPage: req.GetCurrentPage(),
			PageSize:    req.GetPageSize(),
		},
		Roles: convertUserRoleToEntity(req.GetRoles()),
	}

	switch req.GetSortDirection() {
	case userpb.SortDirection_SORT_DIRECTION_ASC:
		request.Paginated.Descending = false
	case userpb.SortDirection_SORT_DIRECTION_DESC:
		request.Paginated.Descending = true
	case userpb.SortDirection_SORT_DIRECTION_UNSPECIFIED:
		request.Paginated.Descending = true
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

func convertUserUpdateFromSuperuserToEntity(req *userpb.UserUpdateFromSuperUserRequest) userservice.UserUpdateFromSuperuserRequest {
	return userservice.UserUpdateFromSuperuserRequest{
		Username:    req.GetUsername(),
		Fullname:    req.GetFullname(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Roles:       convertUserRoleToEntity(req.GetRoles()),
	}
}

func convertCreateUserRequestToEntity(req *userpb.CreateUserRequest) userservice.UserCreateRequest {
	return userservice.UserCreateRequest{
		ID:          "",
		Username:    req.GetUsername(),
		Password:    req.GetPassword(),
		Fullname:    req.GetFullname(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Roles:       convertUserRoleToEntity(req.GetRoles()),
	}
}
