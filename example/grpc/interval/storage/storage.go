package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/storage"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/errorhandler"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/managerauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

func Storage() {
	refreshToken, accessToken := managerauth.GetToken()

	_ = refreshToken

	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2561,
		SSLMode: false,
		UseOtel: true,
	}

	ctxWithValue := context.WithValue(context.Background(), types.AuthorizationKey, "Bearer "+accessToken)
	logger := slog.Default()

	grpcClient, gErr := grpcclient.New(cfg, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
	if gErr != nil {
		panic(gErr)
	}

	storageAd := storage.NewInternalAdapter(grpcClient.Conn())

	getLinkRes, getLinkErr := storageAd.GetLink(ctxWithValue, &storagepb.GetLinkRequest{FileId: "01KCWR8MCTWA9769YFYW7CV779"})
	if getLinkErr != nil {
		errorhandler.HandleGRPCError(getLinkErr, logger)

		return
	}

	fmt.Printf("%+v\n", getLinkRes)

	fmt.Println("getfileInfo")

	getFileInfoRes, getFileInfoErr := storageAd.GetFileInfo(ctxWithValue, &storagepb.GetFileInfoRequest{FileId: "01KCWR8MCTWA9769YFYW7CV779"})
	if getFileInfoErr != nil {
		errorhandler.HandleGRPCError(getFileInfoErr, logger)
	}

	fmt.Printf("%+v\n", getFileInfoRes)
	fmt.Println("is confirmed: ", getFileInfoRes.IsConfirmed)
}
