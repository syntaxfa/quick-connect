package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/protobuf/shared/golang/errdetailspb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
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
	pkResp, pkErr := authAdapter.GetPublicKey(context.Background())
	if pkErr != nil {
		handleGRPCError(pkErr, slog.Default())
	}
	fmt.Println("Public Key:", pkResp)

	fmt.Println("--------------------")

	resp, lErr := authAdapter.Login(context.Background(), &authpb.LoginRequest{
		Username: "alireza",
		Password: "Password",
	})
	if lErr != nil {
		handleGRPCError(lErr, slog.Default())

		return
	}

	fmt.Printf("%+v\n", resp)

	tokenVerifyResp, tvErr := authAdapter.TokenVerify(context.Background(), &authpb.TokenVerifyRequest{Token: resp.RefreshToken})
	if tvErr != nil {
		handleGRPCError(tvErr, slog.Default())
	}

	fmt.Printf("%+v\n", tokenVerifyResp)
}

func handleGRPCError(err error, logger *slog.Logger) {
	st, ok := status.FromError(err)
	if !ok {
		logger.Error("Non-gRPC error occurred", slog.String("error", err.Error()))
		fmt.Printf("An unexpected error occurred: %v\n", err)

		return
	}

	logger.Warn("gRPC request failed",
		slog.String("code", st.Code().String()),
		slog.String("message", st.Message()))
	fmt.Printf("gRPC Error:\n Code: %s\n message: %s\n", st.Code(), st.Message())

	if st.Code() == codes.InvalidArgument {
		fmt.Println(" Details:")
		foundDetails := false

		for _, detail := range st.Details() {
			switch d := detail.(type) {
			case *errdetailspb.BadRequest:
				foundDetails = true
				fmt.Println(" BadRequest Details:")
				for _, violation := range d.GetFieldViolations() {
					fmt.Printf("  - Field: %s, description: %s\n", violation.GetField(), violation.GetDescription())
				}
			default:
				fmt.Printf(" Unknown Detail Type: %T\n", d)
			}
		}
		if !foundDetails {
			fmt.Println("No specific BadRequest details found")
		}
	}
}
