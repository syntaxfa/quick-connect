package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
	"github.com/syntaxfa/quick-connect/types"
)

func (s *Service) GetUserActiveConversation(ctx context.Context, userID types.ID) (ConversationDetailResponse, error) {
	const op = "service.conversation.GetOpenConversation"

	exists, existErr := s.repo.IsUserHaveActiveConversation(ctx, userID)
	if existErr != nil {
		return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	var conversation Conversation

	if !exists {
		id := ulid.Make().String()

		var createErr error
		conversation, createErr = s.repo.CreateActiveConversation(ctx, types.ID(id), userID, ConversationStatusNew)
		if createErr != nil {
			return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(createErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}
	} else {
		var getErr error
		conversation, getErr = s.repo.GetUserActiveConversation(ctx, userID)
		if getErr != nil {
			return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(getErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}
	}

	clientInfo, supportInfo, infoErr := s.getClientAndSupportInfo(ctx, conversation.ClientUserID, conversation.AssignedSupportID)
	if infoErr != nil {
		return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(infoErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return ConversationDetailResponse{
		Conversation: conversation,
		ClientInfo:   clientInfo,
		SupportInfo:  supportInfo,
	}, nil
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

func (s *Service) GetConversationByID(ctx context.Context, conversationID, userID types.ID, userRoles []types.Role) (
	ConversationDetailResponse, error) {
	const op = "service.conversation.GetConversationByID"

	if exists, exErr := s.repo.IsConversationExistByID(ctx, conversationID); exErr != nil {
		return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).
			WithKind(richerror.KindUnexpected), s.logger)
	} else if !exists {
		return ConversationDetailResponse{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).
			WithKind(richerror.KindNotFound)
	}

	var roleOk bool
	for _, role := range userRoles {
		if role == types.RoleSupport || role == types.RoleSuperUser {
			roleOk = true
		}
	}

	if !roleOk {
		ok, checkErr := s.repo.CheckUserInConversation(ctx, userID, conversationID)
		if checkErr != nil {
			return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(checkErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		if !ok {
			return ConversationDetailResponse{}, richerror.New(op).WithMessage(servermsg.MsgConversationNotFound).
				WithKind(richerror.KindNotFound)
		}
	}

	conversation, gErr := s.repo.GetConversationByID(ctx, conversationID)
	if gErr != nil {
		return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	clientInfo, supportInfo, infoErr := s.getClientAndSupportInfo(ctx, conversation.ClientUserID, conversation.AssignedSupportID)
	if infoErr != nil {
		return ConversationDetailResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(infoErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return ConversationDetailResponse{
		Conversation: conversation,
		ClientInfo:   clientInfo,
		SupportInfo:  supportInfo,
	}, nil
}

func (s *Service) getClientAndSupportInfo(ctx context.Context, clientID, supportID types.ID) (ClientInfo, SupportInfo, error) {
	const op = "service.conversation.getClientAndSupportInfo"

	ctxWithAuth, tErr := s.tokenManager.SetTokenInContext(ctx)
	if tErr != nil {
		return ClientInfo{}, SupportInfo{}, errlog.ErrContext(ctxWithAuth, richerror.New(op).WithWrapError(tErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	var clientInfoPB *userinternalpb.UserInfoResponse
	if clientID != "" {
		var gcErr error
		clientInfoPB, gcErr = s.userInternalSvc.UserInfo(ctxWithAuth, &userinternalpb.UserInfoRequest{UserId: string(clientID)})
		if gcErr != nil {
			return ClientInfo{}, SupportInfo{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gcErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}
	}

	var supportInfoPB *userinternalpb.UserInfoResponse
	if supportID != "" {
		var gsErr error
		supportInfoPB, gsErr = s.userInternalSvc.UserInfo(ctxWithAuth, &userinternalpb.UserInfoRequest{UserId: string(supportID)})
		if gsErr != nil {
			return ClientInfo{}, SupportInfo{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gsErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}
	}

	var clientInfo ClientInfo
	var supportInfo SupportInfo

	if clientInfoPB != nil {
		clientInfo = ClientInfo{
			ID:           types.ID(clientInfoPB.GetId()),
			Fullname:     clientInfoPB.GetFullname(),
			PhoneNumber:  clientInfoPB.GetPhoneNumber(),
			Email:        clientInfoPB.GetEmail(),
			Avatar:       clientInfoPB.GetAvatar(),
			LastOnlineAt: clientInfoPB.GetLastOnlineAt().AsTime(),
		}
	}

	if supportInfoPB != nil {
		supportInfo = SupportInfo{
			ID:           types.ID(supportInfoPB.GetId()),
			Fullname:     supportInfoPB.GetFullname(),
			Avatar:       supportInfoPB.GetAvatar(),
			LastOnlineAt: supportInfoPB.GetLastOnlineAt().AsTime(),
		}
	}

	return clientInfo, supportInfo, nil
}
