package micro2

import (
	"context"

	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/service"
	"github.com/syntaxfa/quick-connect/protobuf/example/golang/examplepb"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c Client) GetComment(ctx context.Context, commentID uint64) (service.GetCommentResponse, error) {
	client := examplepb.NewCommentServiceClient(c.conn)

	res, err := client.GetComment(ctx, &examplepb.GetCommentByIDRequest{CommentId: commentID})
	if err != nil {
		return service.GetCommentResponse{}, err
	}

	return service.GetCommentResponse{
		ID:   res.CommentId,
		Body: res.Body,
	}, nil
}
