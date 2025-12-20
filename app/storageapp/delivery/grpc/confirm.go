package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (h InternalHandler) ConfirmFile(ctx context.Context, req *storagepb.ConfirmFileRequest) (*empty.Empty, error) {
	if sErr := h.svc.ConfirmFile(ctx, types.ID(req.GetFileId())); sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return &empty.Empty{}, nil
}
