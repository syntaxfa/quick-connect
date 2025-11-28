package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	GetPublicKey(ctx context.Context, req *empty.Empty, opts ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error)
	TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest, opts ...grpc.CallOption) (*authpb.TokenRefreshResponse, error)
}

type UserService interface {
	UserProfile(ctx context.Context, req *empty.Empty, opts ...grpc.CallOption) (*userpb.User, error)
	UserUpdateFromOwn(ctx context.Context, req *userpb.UserUpdateFromOwnRequest, opts ...grpc.CallOption) (*userpb.User, error)
	UserList(ctx context.Context, req *userpb.UserListRequest, opts ...grpc.CallOption) (*userpb.UserListResponse, error)
	UserDelete(ctx context.Context, req *userpb.UserDeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UserDetail(ctx context.Context, req *userpb.UserDetailRequest, opts ...grpc.CallOption) (*userpb.User, error)
	UserUpdateFromSuperuser(ctx context.Context, req *userpb.UserUpdateFromSuperUserRequest, opts ...grpc.CallOption) (*userpb.User, error)
	CreateUser(ctx context.Context, req *userpb.CreateUserRequest, opts ...grpc.CallOption) (*userpb.User, error)
	UserChangePassword(ctx context.Context, req *userpb.UserChangePasswordRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type ConversationService interface {
	ConversationNewList(ctx context.Context, req *conversationpb.ConversationListRequest,
		opts ...grpc.CallOption) (*conversationpb.ConversationListResponse, error)
	ConversationOwnList(ctx context.Context, req *conversationpb.ConversationListRequest,
		opts ...grpc.CallOption) (*conversationpb.ConversationListResponse, error)
	ConversationDetail(ctx context.Context, req *conversationpb.ConversationDetailRequest,
		opts ...grpc.CallOption) (*conversationpb.ConversationDetailResponse, error)
	ChatHistory(ctx context.Context, req *conversationpb.ChatHistoryRequest,
		opts ...grpc.CallOption) (*conversationpb.ChatHistoryResponse, error)
	OpenConversation(ctx context.Context, req *conversationpb.OpenConversationRequest,
		opts ...grpc.CallOption) (*conversationpb.Conversation, error)
	CloseConversation(ctx context.Context, req *conversationpb.CloseConversationRequest,
		opts ...grpc.CallOption) (*conversationpb.Conversation, error)
}
