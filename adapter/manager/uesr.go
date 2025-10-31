package manager

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"google.golang.org/grpc"
)

// UserAdapter acts as a client adapter for the manager's UserService gRPC service.
type UserAdapter struct {
	client userpb.UserServiceClient
}

// NewUserAdapter creates a new UserAdapter.
// It's generally recommended to pass the connection rather than the client interface
// if the adapter itself doesn't need complex logic, which is the case here.
func NewUserAdapter(conn grpc.ClientConnInterface) *UserAdapter {
	return &UserAdapter{
		client: userpb.NewUserServiceClient(conn),
	}
}

// CreateUser calls the CreateUser RPC on the UserService.
func (ud *UserAdapter) CreateUser(ctx context.Context, req *userpb.CreateUserRequest, opts ...grpc.CallOption) (*userpb.User, error) {
	return ud.client.CreateUser(ctx, req, opts...)
}
