package chat

import (
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/paginate/cursorbased"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertConversationStatusToEntity(pbStatuses []conversationpb.Status) []service.ConversationStatus {
	var statuses []service.ConversationStatus

	for _, status := range pbStatuses {
		switch status {
		case conversationpb.Status_STATUS_NEW:
			statuses = append(statuses, service.ConversationStatusNew)
		case conversationpb.Status_STATUS_OPEN:
			statuses = append(statuses, service.ConversationStatusOpen)
		case conversationpb.Status_STATUS_CLOSED:
			statuses = append(statuses, service.ConversationStatusClosed)
		case conversationpb.Status_STATUS_BOT_HANDLING:
			statuses = append(statuses, service.ConversationStatusBotHandling)
		case conversationpb.Status_STATUS_UNSPECIFIED:
			continue
		}
	}

	return statuses
}

func convertConversationListRequestToEntity(req *conversationpb.ConversationListRequest) service.ListConversationsRequest {
	request := service.ListConversationsRequest{
		AssignedSupportID: types.ID(req.GetAssignedSupportId()),
		Statuses:          convertConversationStatusToEntity(req.GetStatuses()),
		Paginated: paginate.RequestBase{
			CurrentPage: req.GetCurrentPage(),
			PageSize:    req.GetPageSize(),
		},
	}

	switch req.GetSortDirection() {
	case conversationpb.SortDirection_SORT_DIRECTION_ASC:
		request.Paginated.Descending = false
	case conversationpb.SortDirection_SORT_DIRECTION_DESC:
		request.Paginated.Descending = true
	case conversationpb.SortDirection_SORT_DIRECTION_UNSPECIFIED:
		request.Paginated.Descending = true
	default:
		request.Paginated.Descending = true
	}

	return request
}

func convertConversationStatusToPB(status service.ConversationStatus) conversationpb.Status {
	switch status {
	case service.ConversationStatusNew:
		return conversationpb.Status_STATUS_NEW
	case service.ConversationStatusOpen:
		return conversationpb.Status_STATUS_OPEN
	case service.ConversationStatusClosed:
		return conversationpb.Status_STATUS_CLOSED
	case service.ConversationStatusBotHandling:
		return conversationpb.Status_STATUS_BOT_HANDLING
	default:
		return conversationpb.Status_STATUS_UNSPECIFIED
	}
}

func convertConversationToPB(conversation service.Conversation) *conversationpb.Conversation {
	var closedAt *timestamppb.Timestamp

	if conversation.ClosedAt != nil {
		closedAt = timestamppb.New(*conversation.ClosedAt)
	}

	return &conversationpb.Conversation{
		Id:                  string(conversation.ID),
		ClientUserId:        string(conversation.ClientUserID),
		AssignedSupportId:   string(conversation.AssignedSupportID),
		Status:              convertConversationStatusToPB(conversation.Status),
		LastMessageSnippet:  conversation.LastMessageSnippet,
		LastMessageSenderId: string(conversation.LastMessageSenderID),
		CreatedAt:           timestamppb.New(conversation.CreatedAt),
		UpdatedAt:           timestamppb.New(conversation.UpdatedAt),
		ClosedAt:            closedAt,
	}
}

func convertConversationListResponseToPB(resp service.ListConversationsResponse) *conversationpb.ConversationListResponse {
	var cons []*conversationpb.Conversation
	for _, con := range resp.Results {
		cons = append(cons, convertConversationToPB(con))
	}

	return &conversationpb.ConversationListResponse{
		CurrentPage:   resp.Paginate.CurrentPage,
		PageSize:      resp.Paginate.PageSize,
		TotalNumber:   resp.Paginate.TotalNumbers,
		TotalPage:     resp.Paginate.TotalPage,
		Conversations: cons,
	}
}

func convertConversationDetailResponseToPB(resp service.ConversationDetailResponse) *conversationpb.ConversationDetailResponse {
	return &conversationpb.ConversationDetailResponse{
		Conversation: convertConversationToPB(resp.Conversation),
		ClientInfo: &conversationpb.ClientInfo{
			Id:           string(resp.ClientInfo.ID),
			Fullname:     resp.ClientInfo.Fullname,
			PhoneNumber:  resp.ClientInfo.PhoneNumber,
			Email:        resp.ClientInfo.Email,
			Avatar:       resp.ClientInfo.Avatar,
			LastOnlineAt: timestamppb.New(resp.ClientInfo.LastOnlineAt),
		},
		SupportInfo: &conversationpb.SupportInfo{
			Id:           string(resp.SupportInfo.ID),
			Fullname:     resp.SupportInfo.Fullname,
			Avatar:       resp.SupportInfo.Avatar,
			LastOnlineAt: timestamppb.New(resp.SupportInfo.LastOnlineAt),
		},
	}
}

func convertChatHistoryRequestToEntity(req *conversationpb.ChatHistoryRequest) service.ChatHistoryRequest {
	return service.ChatHistoryRequest{
		UserID:         "",
		UserRoles:      nil,
		ConversationID: types.ID(req.GetConversationId()),
		Pagination: cursorbased.Request{
			Cursor: types.ID(req.GetCursor()),
			Limit:  int(req.GetLimit()),
		},
	}
}

func convertMessageTypeToPB(messageType service.MessageType) conversationpb.MessageType {
	switch messageType {
	case service.MessageTypeText:
		return conversationpb.MessageType_TYPE_TEXT
	case service.MessageTypeMedia:
		return conversationpb.MessageType_TYPE_MEDIA
	case service.MessageTypeSystem:
		return conversationpb.MessageType_TYPE_SYSTEM
	default:
		return conversationpb.MessageType_TYPE_UNSPECIFIED
	}
}

func convertMessageToPB(message service.Message) *conversationpb.Message {
	var readAt *timestamppb.Timestamp

	if message.ReadAt != nil {
		readAt = timestamppb.New(*message.ReadAt)
	}

	return &conversationpb.Message{
		Id:                 string(message.ID),
		ConversationId:     string(message.ConversationID),
		SenderId:           string(message.SenderID),
		MessageType:        convertMessageTypeToPB(message.MessageType),
		Content:            message.Content,
		Metadata:           message.Metadata,
		RepliedToMessageId: string(message.RepliedToMessageID),
		CreatedAt:          timestamppb.New(message.CreatedAt),
		ReadAt:             readAt,
	}
}

func convertChatHistoryResponseToPB(resp service.ChatHistoryResponse) *conversationpb.ChatHistoryResponse {
	var results []*conversationpb.Message

	for _, msg := range resp.Results {
		results = append(results, convertMessageToPB(msg))
	}

	return &conversationpb.ChatHistoryResponse{
		Results:    results,
		NextCursor: string(resp.Paginate.NextCursor),
		HasMore:    resp.Paginate.HasMore,
	}
}
