package storage

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type InternalAdapter struct {
	client storagepb.StorageInternalServiceClient
}

func NewInternalAdapter(conn grpc.ClientConnInterface) *InternalAdapter {
	return &InternalAdapter{
		client: storagepb.NewStorageInternalServiceClient(conn),
	}
}

func (id *InternalAdapter) GetLink(ctx context.Context, req *storagepb.GetLinkRequest,
	opts ...grpc.CallOption) (*storagepb.GetLinkResponse, error) {
	return id.client.GetLink(ctx, req, opts...)
}

func (id *InternalAdapter) GetFileInfo(ctx context.Context, req *storagepb.GetFileInfoRequest,
	opts ...grpc.CallOption) (*storagepb.File, error) {
	return id.client.GetFileInfo(ctx, req, opts...)
}

func (id *InternalAdapter) ConfirmFile(ctx context.Context, req *storagepb.ConfirmFileRequest,
	opts ...grpc.CallOption) (*empty.Empty, error) {
	return id.client.ConfirmFile(ctx, req, opts...)
}
