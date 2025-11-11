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
