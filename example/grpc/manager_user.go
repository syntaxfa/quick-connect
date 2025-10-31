package main

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

func main() {
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

	ctxWithValue := context.WithValue(context.Background(), types.AuthorizationKey, "Bearer "+resp.AccessToken)

	createUserResp, createErr := userAdapter.CreateUser(ctxWithValue, &userpb.CreateUserRequest{
		Username:    "ayda",
		Password:    "Password",
		Fullname:    "ayda jon",
		Email:       "ayda@gmail.com",
		PhoneNumber: "09307225656",
		Roles:       []userpb.Role{userpb.Role_ROLE_SUPPORT, userpb.Role_ROLE_STORY},
	})
	if createErr != nil {
		errorhandler.HandleGRPCError(createErr, slog.Default())
	}

	fmt.Printf("%+v\n", createUserResp)

	fmt.Println("-----------------")
	fmt.Println("User Detail:")

	ctxWithValue = context.WithValue(context.Background(), types.AuthorizationKey, "Bearer "+resp.AccessToken)

	userDetailResp, userDetailErr := userAdapter.UserDetail(ctxWithValue, &userpb.UserDetailRequest{UserId: "01K8QVTKGM9XRKV29T7BAPAK9J"})
	if userDetailErr != nil {
		errorhandler.HandleGRPCError(createErr, slog.Default())
	}

	fmt.Printf("%+v\n", userDetailResp)

	fmt.Println("-----------------")
}
