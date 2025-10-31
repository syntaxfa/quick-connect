package manager

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
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

// UserDetail calls the UserDetail PRC on the UserService.
func (ud *UserAdapter) UserDetail(ctx context.Context, req *userpb.UserDetailRequest, opts ...grpc.CallOption) (*userpb.User, error) {
	return ud.client.UserDetail(ctx, req, opts...)
}

// UserDelete calls the UserDelete PRC on the UserService.
func (ud *UserAdapter) UserDelete(ctx context.Context, req *userpb.UserDeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return ud.client.UserDelete(ctx, req, opts...)
}

// UserUpdateFromSuperuser calls the UserUpdateFromSuperuser PRC on the UserService.
func (ud *UserAdapter) UserUpdateFromSuperuser(ctx context.Context, req *userpb.UserUpdateFromSuperUserRequest, opts ...grpc.CallOption) (*userpb.User, error) {
	return ud.client.UserUpdateFromSuperuser(ctx, req, opts...)
}

// UserList calls the UserList PRC on the UserService.
func (ud *UserAdapter) UserList(ctx context.Context, req *userpb.UserListRequest, opts ...grpc.CallOption) (*userpb.UserListResponse, error) {
	return ud.client.UserList(ctx, req, opts...)
}
