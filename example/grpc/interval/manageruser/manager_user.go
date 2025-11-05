package manageruser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/errorhandler"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

func ManagerUser() {
	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2541,
		SSLMode: false,
		UseOtel: true,
	}

	grpcClient, gErr := grpcclient.New(cfg, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
	if gErr != nil {
		panic(gErr)
	}

	authAdapter := manager.NewAuthAdapter(grpcClient.Conn())

	resp, lErr := authAdapter.Login(context.Background(), &authpb.LoginRequest{
		Username: "alireza",
		Password: "Password",
	})
	if lErr != nil {
		errorhandler.HandleGRPCError(lErr, slog.Default())

		return
	}

	_ = resp

	userAdapter := manager.NewUserAdapter(grpcClient.Conn())

	ctxWithValue := context.WithValue(context.Background(), types.AuthorizationKey, "Bearer "+resp.GetAccessToken())

	//createUserResp, createErr := userAdapter.CreateUser(ctxWithValue, &userpb.CreateUserRequest{
	//	Username:    "ayda",
	//	Password:    "Password",
	//	Fullname:    "ayda jon",
	//	Email:       "ayda@gmail.com",
	//	PhoneNumber: "09307225656",
	//	Roles:       []userpb.Role{userpb.Role_ROLE_SUPPORT, userpb.Role_ROLE_STORY},
	//})
	//if createErr != nil {
	//	errorhandler.HandleGRPCError(createErr, slog.Default())
	//
	//	return
	//}
	//
	//fmt.Printf("%+v\n", createUserResp)

	fmt.Println("-----------------")
	fmt.Println("User Detail:")

	userDetailResp, userDetailErr := userAdapter.UserDetail(ctxWithValue, &userpb.UserDetailRequest{UserId: "01K8XSS8B8XYGM1Y5DYWRXG0S3"})
	if userDetailErr != nil {
		errorhandler.HandleGRPCError(userDetailErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", userDetailResp)

	//fmt.Println("-----------------")
	//fmt.Println("User Delete:")
	//
	//if _, userDeleteErr := userAdapter.UserDelete(ctxWithValue, &userpb.UserDeleteRequest{UserId: "01K8QVTKGM9XRKV29T7BAPAK9J"}); userDeleteErr != nil {
	//	errorhandler.HandleGRPCError(userDeleteErr, slog.Default())
	//
	//	return
	//}
	//
	fmt.Println("-----------------")
	fmt.Println("User Update:")

	userUpResp, UserUpErr := userAdapter.UserUpdateFromSuperuser(ctxWithValue, &userpb.UserUpdateFromSuperUserRequest{
		UserId:      "01K8XSS8B8XYGM1Y5DYWRXG0S3",
		Username:    userDetailResp.GetUsername(),
		Fullname:    "ayda family",
		Email:       "aydafamily@gmail.com",
		PhoneNumber: "09119111111",
		Roles:       []userpb.Role{userpb.Role_ROLE_SUPPORT, userpb.Role_ROLE_NOTIFICATION},
	})
	if UserUpErr != nil {
		errorhandler.HandleGRPCError(UserUpErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", userUpResp)

	fmt.Println("-----------------")
	fmt.Println("User List:")

	userListResp, userListErr := userAdapter.UserList(ctxWithValue, &userpb.UserListRequest{
		CurrentPage:   1,
		PageSize:      10,
		SortDirection: 1,
		Username:      "",
	})
	if userListErr != nil {
		errorhandler.HandleGRPCError(userListErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", userListResp)

	fmt.Println("-----------------")
	fmt.Println("change password")
	if _, userChangeErr := userAdapter.UserChangePassword(ctxWithValue, &userpb.UserChangePasswordRequest{
		OldPassword: "Password",
		NewPassword: "Berlin11228",
	}); userChangeErr != nil {
		errorhandler.HandleGRPCError(userChangeErr, slog.Default())
	}
	fmt.Println("password changed")
}
