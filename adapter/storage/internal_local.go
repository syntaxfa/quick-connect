package storage

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
)

type InternalLocalAdapter struct {
	svc    *service.Service
	t      *translation.Translate
	logger *slog.Logger
}

func NewInternalLocalAdapter(svc *service.Service) *InternalLocalAdapter {
	return &InternalLocalAdapter{
		svc: svc,
	}
}

func (idl *InternalLocalAdapter) GetLink(ctx context.Context, req *storagepb.GetLinkRequest,
	_ ...grpc.CallOption) (*storagepb.GetLinkResponse, error) {
	resp, sErr := idl.svc.GetLink(ctx, types.ID(req.GetFileId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, idl.t, idl.logger)
	}

	return &storagepb.GetLinkResponse{Url: resp}, nil
}

func (idl *InternalLocalAdapter) GetFileInfo(ctx context.Context, req *storagepb.GetFileInfoRequest,
	_ ...grpc.CallOption) (*storagepb.File, error) {
	resp, sErr := idl.svc.GetFileInfo(ctx, types.ID(req.GetFileId()))
	if sErr != nil {
		return nil, servermsg.GRPCMsg(sErr, idl.t, idl.logger)
	}

	return convertFileToPB(resp), nil
}
