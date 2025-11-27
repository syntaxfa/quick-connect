package manager

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
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
