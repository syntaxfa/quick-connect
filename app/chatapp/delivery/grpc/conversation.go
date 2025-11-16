package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handler) ConversationNewList(ctx context.Context, req *conversationpb.ConversationListRequest) (
	*conversationpb.ConversationListResponse, error) {
	userClaims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	request := convertConversationListRequestToEntity(req)
	request.AssignedSupportID = userClaims.UserID

	resp, sErr := h.chatSvc.ListConversations(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertConversationListResponseToPB(resp), nil
}
