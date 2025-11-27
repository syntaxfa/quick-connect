package manager

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

// AuthAdapter acts as a client adapter for the manager's AuthService gRPC service.
type AuthAdapter struct {
	client authpb.AuthServiceClient
}

// NewAuthAdapter creates a new AuthAdapter.
// It's generally recommended to pass the connection rather than the client interface
// if the adapter itself doesn't need complex logic, which is the case here.
func NewAuthAdapter(conn grpc.ClientConnInterface) *AuthAdapter {
	return &AuthAdapter{
		client: authpb.NewAuthServiceClient(conn),
	}
}

// GetPublicKey retrieves the token public key from the AuthService.
func (ad *AuthAdapter) GetPublicKey(ctx context.Context, req *empty.Empty, opts ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error) {
	return ad.client.GetPublicKey(ctx, req, opts...)
}

// Login calls the Login RPC on the AuthService.
func (ad *AuthAdapter) Login(ctx context.Context, req *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	return ad.client.Login(ctx, req, opts...)
}

// TokenVerify calls the TokenVerify PRC on the AuthService.
func (ad *AuthAdapter) TokenVerify(ctx context.Context, req *authpb.TokenVerifyRequest,
	opts ...grpc.CallOption) (*authpb.TokenVerifyResponse, error) {
	return ad.client.TokenVerify(ctx, req, opts...)
}

// TokenRefresh calls the TokenRefresh RPC on the AuthService.
func (ad *AuthAdapter) TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest,
	opts ...grpc.CallOption) (*authpb.TokenRefreshResponse, error) {
	return ad.client.TokenRefresh(ctx, req, opts...)
}
