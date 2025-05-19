package grpc

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/service"
	"github.com/syntaxfa/quick-connect/protobuf/example/golang/examplepb"
)

type Handler struct {
	examplepb.UnimplementedCommentServiceServer
	svc service.Service
}

func NewHandler(svc service.Service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) GetComment(ctx context.Context, req *examplepb.GetCommentByIDRequest) (*examplepb.GetCommentResponse, error) {
	resp, err := h.svc.GetCommentByID(ctx, req.CommentId)
	if err != nil {
		return nil, err
	}

	fmt.Println("everything is ok")
	return &examplepb.GetCommentResponse{CommentId: resp.ID, Body: resp.Body}, nil
}
