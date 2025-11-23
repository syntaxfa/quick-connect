package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handler) ConversationNewList(ctx context.Context, req *conversationpb.ConversationListRequest) (
	*conversationpb.ConversationListResponse, error) {
	request := convertConversationListRequestToEntity(req)
	request.Statuses = []service.ConversationStatus{service.ConversationStatusNew}

	resp, sErr := h.chatSvc.ListConversations(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertConversationListResponseToPB(resp), nil
}

func (h Handler) ConversationOwnList(ctx context.Context, req *conversationpb.ConversationListRequest) (
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

func (h Handler) OpenConversation(ctx context.Context, req *conversationpb.OpenConversationRequest) (*conversationpb.Conversation, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := h.chatSvc.OpenConversation(ctx, types.ID(req.GetConversationId()), claims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertConversationToPB(resp), nil
}

func (h Handler) CloseConversation(ctx context.Context, req *conversationpb.CloseConversationRequest) (
	*conversationpb.Conversation, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := h.chatSvc.CloseConversation(ctx, types.ID(req.GetConversationId()), claims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertConversationToPB(resp), nil
}

func (h Handler) ConversationDetail(ctx context.Context, req *conversationpb.ConversationDetailRequest) (
	*conversationpb.ConversationDetailResponse, error) {
	claims, ucErr := grpcauth.ExtractUserClaimsFromContext(ctx)
	if ucErr != nil {
		return nil, status.Error(codes.Unauthenticated, ucErr.Error())
	}

	resp, sErr := h.chatSvc.GetConversationByID(ctx, types.ID(req.GetConversationId()), claims.UserID, claims.Roles)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertConversationDetailResponseToPB(resp), nil
}
