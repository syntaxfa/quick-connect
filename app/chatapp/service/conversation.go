package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s *Service) GetUserActiveConversation(ctx context.Context, userID types.ID) (Conversation, error) {
	const op = "service.conversation.GetOpenConversation"

	exists, existErr := s.repo.IsUserHaveActiveConversation(ctx, userID)
	if existErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		id := ulid.Make().String()

		conversation, createErr := s.repo.CreateActiveConversation(ctx, types.ID(id), userID, ConversationStatusNew)
		if createErr != nil {
			return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(createErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		return conversation, nil
	}

	conversation, getErr := s.repo.GetUserActiveConversation(ctx, userID)
	if getErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(getErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return conversation, nil
}

// ListConversations handles the business logic for listing conversations.
func (s *Service) ListConversations(ctx context.Context, req ListConversationsRequest) (ListConversationsResponse, error) {
	const op = "service.conversation.ListConversations"

	if bErr := req.Paginated.BasicValidation(); bErr != nil {
		return ListConversationsResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest).
			WithWrapError(bErr)
	}

	if vErr := s.vld.ValidateListConversationsRequest(req); vErr != nil {
		return ListConversationsResponse{}, vErr
	}

	convos, paginateRes, rErr := s.repo.GetConversationList(ctx, req.Paginated, req.AssignedSupportID, req.Statuses)
	if rErr != nil {
		return ListConversationsResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(rErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return ListConversationsResponse{
		Results:  convos,
		Paginate: paginateRes,
	}, nil
}

func (s *Service) OpenConversation(ctx context.Context, conversationID, supportID types.ID) (Conversation, error) {
	const op = "service.conversation.OpenConversation"

	exists, isExErr := s.repo.IsConversationExistByID(ctx, conversationID)
	if isExErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(isExErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return Conversation{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).WithKind(richerror.KindNotFound)
	}

	conversation, gErr := s.repo.GetConversationByID(ctx, conversationID)
	if gErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if conversation.Status != ConversationStatusNew {
		return Conversation{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).WithKind(richerror.KindNotFound)
	}

	if assignErr := s.repo.AssignConversation(ctx, conversationID, supportID); assignErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(assignErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	conversation.AssignedSupportID = supportID
	conversation.Status = ConversationStatusOpen

	return conversation, nil
}

func (s *Service) CloseConversation(ctx context.Context, conversationID, supportID types.ID) (Conversation, error) {
	const op = "service.conversation.CloseConversation"

	exists, isExErr := s.repo.IsConversationExistByID(ctx, conversationID)
	if isExErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(isExErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return Conversation{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).WithKind(richerror.KindNotFound)
	}

	conversation, gErr := s.repo.GetConversationByID(ctx, conversationID)
	if gErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if conversation.AssignedSupportID != supportID {
		return Conversation{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).WithKind(richerror.KindNotFound)
	}

	if conversation.Status == ConversationStatusClosed {
		return Conversation{}, richerror.New(op).WithMessage(servermsg.MsgConversationAlreadyClosed).WithKind(richerror.KindConflict)
	}

	if closeErr := s.repo.CloseConversation(ctx, conversationID); closeErr != nil {
		return Conversation{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(closeErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	conversation.Status = ConversationStatusClosed

	return conversation, nil
}
