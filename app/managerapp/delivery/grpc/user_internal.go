package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
	"github.com/syntaxfa/quick-connect/types"
)

func (h HandlerInternal) UserInfo(ctx context.Context, req *userinternalpb.UserInfoRequest) (*userinternalpb.UserInfoResponse, error) {
	resp, sErr := h.userSvc.UserInfo(ctx, types.ID(req.GetUserId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertUserInfoToPB(resp), nil
}
