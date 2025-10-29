package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) TokenVerify(_ context.Context, req *authpb.TokenVerifyRequest) (*authpb.TokenVerifyResponse, error) {
	userClaims, sErr := h.tokenSvc.ValidateToken(req.Token)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	var roles []string
	for _, role := range userClaims.Roles {
		roles = append(roles, string(role))
	}

	return &authpb.TokenVerifyResponse{
		UserId:    string(userClaims.UserID),
		Roles:     roles,
		TokenType: string(userClaims.TokenType),
	}, nil
}
