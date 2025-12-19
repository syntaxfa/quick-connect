package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
)

func (h InternalHandler) GetLink(ctx context.Context, req *storagepb.GetLinkRequest) (*storagepb.GetLinkResponse, error) {
	resp, sErr := h.svc.GetLink(ctx, types.ID(req.GetFileId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return &storagepb.GetLinkResponse{Url: resp}, nil
}
