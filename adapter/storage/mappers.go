package storage

import (
	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertFileToPB(file service.File) *storagepb.File {
	var deletedAt *timestamppb.Timestamp
	if file.DeletedAt != nil {
		deletedAt = timestamppb.New(*file.DeletedAt)
	}

	return &storagepb.File{
		Id:          string(file.ID),
		UploaderId:  string(file.ID),
		Name:        file.Name,
		Key:         file.Key,
		MimeType:    file.MimeType,
		Size:        file.Size,
		Driver:      string(file.Driver),
		Bucket:      file.Bucket,
		IsPublic:    file.IsPublic,
		IsConfirmed: file.IsConfirmed,
		CreatedAt:   timestamppb.New(file.CreatedAt),
		UpdatedAt:   timestamppb.New(file.UpdatedAt),
		DeletedAt:   deletedAt,
	}
}
