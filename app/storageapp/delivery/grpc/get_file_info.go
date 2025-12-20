package grpc

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
)

func (h InternalHandler) GetFileInfo(ctx context.Context, req *storagepb.GetFileInfoRequest) (*storagepb.File, error) {
	resp, sErr := h.svc.GetFileInfo(ctx, types.ID(req.GetFileId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, h.t, h.logger)
	}

	return convertFileToPB(resp), nil
}
