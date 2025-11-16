package grpc

import "github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"

type Handler struct {
	conversationpb.UnimplementedConversationServiceServer
}

func NewHandler() Handler {
	return Handler{}
}
