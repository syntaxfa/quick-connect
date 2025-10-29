package manager

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc"
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
func (tc AuthAdapter) GetPublicKey(ctx context.Context, opts ...grpc.CallOption) (*authpb.GetPublicKeyResponse, error) {
	return tc.client.GetPublicKey(ctx, &empty.Empty{}, opts...)
}

// Login calls the Login RPC on the AuthService.
func (tc AuthAdapter) Login(ctx context.Context, req *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	return tc.client.Login(ctx, req, opts...)
}

// TokenVerify calls the TokenVerify PRC on the AuthService.
func (tc AuthAdapter) TokenVerify(ctx context.Context, req *authpb.TokenVerifyRequest, opts ...grpc.CallOption) (*authpb.TokenVerifyResponse, error) {
	return tc.client.TokenVerify(ctx, req, opts...)
}

// TokenRefresh calls the TokenRefresh RPC on the AuthService.
func (tc AuthAdapter) TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest, opts ...grpc.CallOption) (*authpb.TokenRefreshResponse, error) {
	return tc.client.TokenRefresh(ctx, req, opts...)
}
