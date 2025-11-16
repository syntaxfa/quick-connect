package chatconversation

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/pkg/grpcauth"
	"google.golang.org/grpc"
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/chat"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/errorhandler"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/managerauth"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"github.com/syntaxfa/quick-connect/types"
)

func ChatConversation() {
	refreshToken, accessToken := managerauth.GetToken()

	_ = refreshToken

	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2551,
		SSLMode: false,
		UseOtel: true,
	}

	ctxWithValue := context.WithValue(context.Background(), types.AuthorizationKey, "Bearer "+accessToken)
	logger := slog.Default()

	grpcClient, gErr := grpcclient.New(cfg, grpc.WithUnaryInterceptor(grpcauth.AuthClientInterceptor))
	if gErr != nil {
		panic(gErr)
	}

	conAd := chat.NewConversationAdapter(grpcClient.Conn())

	conNewList, conNewListErr := conAd.ConversationNewList(ctxWithValue, &conversationpb.ConversationListRequest{
		CurrentPage:   0,
		PageSize:      10,
		SortDirection: 1,
		Statuses:      nil,
	})
	if conNewListErr != nil {
		errorhandler.HandleGRPCError(conNewListErr, logger)

		return
	}

	fmt.Printf("%+v\n", conNewList)

	conOwnList, conOwnListErr := conAd.ConversationOwnList(ctxWithValue, &conversationpb.ConversationListRequest{
		CurrentPage:   0,
		PageSize:      10,
		SortDirection: 1,
		Statuses:      []conversationpb.Status{conversationpb.Status_STATUS_OPEN, conversationpb.Status_STATUS_CLOSED},
	})
	if conOwnListErr != nil {
		errorhandler.HandleGRPCError(conOwnListErr, logger)

		return
	}

	fmt.Printf("%+v\n", conOwnList)
}
