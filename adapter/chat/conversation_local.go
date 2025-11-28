package chat

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/rolemanager"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConversationLocalAdapter struct {
	chatSvc      *service.Service
	t            *translation.Translate
	logger       *slog.Logger
	roleManager  *rolemanager.RoleManager
	jwtValidator *jwtvalidator.Validator
}

func NewConversationLocalAdapter(chatSvc *service.Service, t *translation.Translate, logger *slog.Logger,
	roleManager *rolemanager.RoleManager, jwtValidator *jwtvalidator.Validator) *ConversationLocalAdapter {
	return &ConversationLocalAdapter{
		chatSvc:      chatSvc,
		t:            t,
		logger:       logger,
		roleManager:  roleManager,
		jwtValidator: jwtValidator,
	}
}

func (cal *ConversationLocalAdapter) ConversationNewList(ctx context.Context, req *conversationpb.ConversationListRequest,
	_ ...grpc.CallOption) (*conversationpb.ConversationListResponse, error) {
	_, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/ConversationNewList")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	request := convertConversationListRequestToEntity(req)
	request.Statuses = []service.ConversationStatus{service.ConversationStatusNew}

	resp, sErr := cal.chatSvc.ListConversations(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertConversationListResponseToPB(resp), nil
}

func (cal *ConversationLocalAdapter) ConversationOwnList(ctx context.Context, req *conversationpb.ConversationListRequest,
	_ ...grpc.CallOption) (*conversationpb.ConversationListResponse, error) {
	claims, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/ConversationOwnList")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	request := convertConversationListRequestToEntity(req)
	request.AssignedSupportID = claims.UserID

	resp, sErr := cal.chatSvc.ListConversations(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertConversationListResponseToPB(resp), nil
}

func (cal *ConversationLocalAdapter) ConversationDetail(ctx context.Context, req *conversationpb.ConversationDetailRequest,
	_ ...grpc.CallOption) (*conversationpb.ConversationDetailResponse, error) {
	claims, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/ConversationDetail")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	resp, sErr := cal.chatSvc.GetConversationByID(ctx, types.ID(req.GetConversationId()), claims.UserID, claims.Roles)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertConversationDetailResponseToPB(resp), nil
}

func (cal *ConversationLocalAdapter) ChatHistory(ctx context.Context, req *conversationpb.ChatHistoryRequest,
	_ ...grpc.CallOption) (*conversationpb.ChatHistoryResponse, error) {
	claims, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/ChatHistory")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	request := convertChatHistoryRequestToEntity(req)
	request.UserID = claims.UserID
	request.UserRoles = claims.Roles

	resp, sErr := cal.chatSvc.ChatHistory(ctx, request)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertChatHistoryResponseToPB(resp), nil
}

func (cal *ConversationLocalAdapter) OpenConversation(ctx context.Context, req *conversationpb.OpenConversationRequest,
	_ ...grpc.CallOption) (*conversationpb.Conversation, error) {
	claims, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/OpenConversation")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	resp, sErr := cal.chatSvc.OpenConversation(ctx, types.ID(req.GetConversationId()), claims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertConversationToPB(resp), nil
}

func (cal *ConversationLocalAdapter) CloseConversation(ctx context.Context, req *conversationpb.CloseConversationRequest,
	_ ...grpc.CallOption) (*conversationpb.Conversation, error) {
	claims, pErr := grpcauth.Protect(ctx, cal.roleManager, cal.jwtValidator, "/chat.ConversationService/CloseConversation")
	if pErr != nil {
		return nil, status.Error(codes.Unauthenticated, pErr.Error())
	}

	resp, sErr := cal.chatSvc.CloseConversation(ctx, types.ID(req.GetConversationId()), claims.UserID)
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, cal.t, cal.logger)
	}

	return convertConversationToPB(resp), nil
}
