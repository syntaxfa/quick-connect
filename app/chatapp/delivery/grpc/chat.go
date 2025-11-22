package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handler) ChatHistory(ctx context.Context, req *conversationpb.ChatHistoryRequest) (*conversationpb.ChatHistoryResponse, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	request := convertChatHistoryRequestToEntity(req)
	request.UserID = claims.UserID
	request.UserRoles = claims.Roles

	resp, sErr := h.chatSvc.ChatHistory(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertChatHistoryResponseToPB(resp), nil
}
