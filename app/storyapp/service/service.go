package service

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/tokenmanager"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type Repository interface {
	SaveStory(ctx context.Context, req AddStoryRequest) (Story, error)
}

type StorageService interface {
	GetFileInfo(ctx context.Context, req *storagepb.GetFileInfoRequest, opts ...grpc.CallOption) (*storagepb.File, error)
	ConfirmFile(ctx context.Context, req *storagepb.ConfirmFileRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type Service struct {
	repo         Repository
	vld          Validate
	storageSvc   StorageService
	tokenManager *tokenmanager.TokenManager
	logger       *slog.Logger
}

func New(repo Repository, vld Validate, storageService StorageService, tokenManager *tokenmanager.TokenManager,
	logger *slog.Logger) Service {
	return Service{
		repo:         repo,
		vld:          vld,
		storageSvc:   storageService,
		tokenManager: tokenManager,
		logger:       logger,
	}
}
