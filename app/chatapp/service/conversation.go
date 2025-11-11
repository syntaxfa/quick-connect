package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
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
