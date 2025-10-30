package grpcauth

import (
	"context"
	"fmt"

	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("AuthClientInterceptor working ...")

	token, ok := ctx.Value(types.AuthorizationKey).(string)
	if !ok {
		fmt.Println("not set authorization header")
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, types.AuthorizationKey, token)

	return invoker(ctx, method, req, reply, cc, opts...)
}
