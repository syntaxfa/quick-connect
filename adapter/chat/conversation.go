package chat

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"google.golang.org/grpc"
)

// ConversationAdapter acts as a client adapter for the chat's ConversationService gRPC service.
type ConversationAdapter struct {
	client conversationpb.ConversationServiceClient
}

func NewConversationAdapter(conn grpc.ClientConnInterface) *ConversationAdapter {
	return &ConversationAdapter{
		client: conversationpb.NewConversationServiceClient(conn),
	}
}

func (ca *ConversationAdapter) ConversationNewList(ctx context.Context, req *conversationpb.ConversationListRequest,
	opts ...grpc.CallOption) (*conversationpb.ConversationListResponse, error) {
	return ca.client.ConversationNewList(ctx, req, opts...)
}

func (ca *ConversationAdapter) ConversationOwnList(ctx context.Context, req *conversationpb.ConversationListRequest,
	opts ...grpc.CallOption) (*conversationpb.ConversationListResponse, error) {
	return ca.client.ConversationOwnList(ctx, req, opts...)
}

func (ca *ConversationAdapter) ChatHistory(ctx context.Context, req *conversationpb.ChatHistoryRequest,
	opts ...grpc.CallOption) (*conversationpb.ChatHistoryResponse, error) {
	return ca.client.ChatHistory(ctx, req, opts...)
}
