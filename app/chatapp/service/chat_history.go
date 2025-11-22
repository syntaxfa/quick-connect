package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s *Service) ChatHistory(ctx context.Context, req ChatHistoryRequest) (ChatHistoryResponse, error) {
	const op = "service.chat_history.ChatHistory"

	if bvErr := req.Pagination.BasicValidation(); bvErr != nil {
		return ChatHistoryResponse{}, richerror.New(op).WithWrapError(bvErr).WithMessage(servermsg.MsgInvalidInput).
			WithKind(richerror.KindBadRequest)
	}

	if userInConv, checkErr := s.checkUserParticipantInConversation(ctx, req.ConversationID, req.UserID, req.UserRoles); checkErr != nil {
		return ChatHistoryResponse{}, checkErr
	} else if !userInConv {
		return ChatHistoryResponse{}, errlog.ErrContext(ctx, richerror.New(op).
			WithMessage(servermsg.MsgYouNotParticipantInThisConversation).WithKind(richerror.KindForbidden).
			WithMeta(map[string]interface{}{"user_id": req.UserID, "conversation_id": req.ConversationID}), s.logger)
	}

	resp, gErr := s.repo.GetChatHistory(ctx, req.ConversationID, req.Pagination)
	if gErr != nil {
		return ChatHistoryResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return resp, nil
}

func (s *Service) checkUserParticipantInConversation(ctx context.Context, conversationID, userID types.ID,
	userRoles []types.Role) (bool, error) {
	const op = "service.chat_history.checkUserParticipantInConversation"

	for _, role := range userRoles {
		if role == types.RoleSupport || role == types.RoleSuperUser {
			return true, nil
		}
	}

	isParticipant, ispErr := s.repo.CheckUserInConversation(ctx, userID, conversationID)
	if ispErr != nil {
		return false, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(ispErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !isParticipant {
		return false, nil
	}

	return true, nil
}
