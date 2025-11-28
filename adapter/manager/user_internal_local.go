package manager

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

// UserInternalLocalAdapter acts as a client local adapter with func calls.
type UserInternalLocalAdapter struct {
	userSvc *userservice.Service
}

func NewUserInternalLocalAdapter(userSvc *userservice.Service) *UserInternalLocalAdapter {
	return &UserInternalLocalAdapter{
		userSvc: userSvc,
	}
}

func (uil *UserInternalLocalAdapter) UserInfo(ctx context.Context, req *userinternalpb.UserInfoRequest, _ ...grpc.CallOption) (
	*userinternalpb.UserInfoResponse, error) {
	resp, sErr := uil.userSvc.UserInfo(ctx, types.ID(req.GetUserId()))
	if sErr != nil {
		return nil, sErr
	}

	return convertUserInfoToPB(resp), nil
}
