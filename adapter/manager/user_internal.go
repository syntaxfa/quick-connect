package manager

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
	"google.golang.org/grpc"
)

// UserInternalAdapter acts as a client adapter for the manager's UserService internal gRPC service.
type UserInternalAdapter struct {
	client userinternalpb.UserInternalServiceClient
}

func NewUserInternalAdapter(conn grpc.ClientConnInterface) *UserInternalAdapter {
	return &UserInternalAdapter{
		client: userinternalpb.NewUserInternalServiceClient(conn),
	}
}

// UserInfo calls the UserInfo on the UserService.
func (ui *UserInternalAdapter) UserInfo(ctx context.Context, req *userinternalpb.UserInfoRequest,
	opts ...grpc.CallOption) (*userinternalpb.UserInfoResponse, error) {
	return ui.client.UserInfo(ctx, req, opts...)
}
