package grpcauth

import (
	"context"

	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	token, ok := ctx.Value(types.AuthorizationKey).(string)
	if !ok {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, string(types.AuthorizationKey), token)

	return invoker(ctx, method, req, reply, cc, opts...)
}
