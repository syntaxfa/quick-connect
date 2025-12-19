package storage

import (
	"context"

	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"google.golang.org/grpc"
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
