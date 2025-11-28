package managerauth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/errorhandler"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func GetToken() (refreshToken, accessToken string) {
	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2541,
		SSLMode: false,
		UseOtel: true,
	}

	grpcClient, gErr := grpcclient.New(cfg)
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

	return resp.GetRefreshToken(), resp.GetAccessToken()
}

func ManagerAuth() {
	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2541,
		SSLMode: false,
		UseOtel: true,
	}

	grpcClient, gErr := grpcclient.New(cfg)
	if gErr != nil {
		panic(gErr)
	}

	authAdapter := manager.NewAuthAdapter(grpcClient.Conn())
	pkResp, pkErr := authAdapter.GetPublicKey(context.Background(), nil)
	if pkErr != nil {
		errorhandler.HandleGRPCError(pkErr, slog.Default())
	}
	fmt.Println("Public Key:", pkResp)

	fmt.Println("--------------------")

	resp, lErr := authAdapter.Login(context.Background(), &authpb.LoginRequest{
		Username: "alireza",
		Password: "Password",
	})
	if lErr != nil {
		errorhandler.HandleGRPCError(lErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", resp)

	tokenVerifyResp, tvErr := authAdapter.TokenVerify(context.Background(), &authpb.TokenVerifyRequest{Token: resp.RefreshToken})
	if tvErr != nil {
		errorhandler.HandleGRPCError(tvErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", tokenVerifyResp)

	fmt.Println("---------------------------------")
	fmt.Println("Token refresh:")

	tokenRefreshResp, tfErr := authAdapter.TokenRefresh(context.Background(), &authpb.TokenRefreshRequest{RefreshToken: resp.RefreshToken})
	if tfErr != nil {
		errorhandler.HandleGRPCError(tfErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", tokenRefreshResp)
}
