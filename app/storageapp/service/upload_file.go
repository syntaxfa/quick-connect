package service

import "context"

func (s Service) Upload(_ context.Context, _ UploadRequest) (File, error) {
	return File{}, nil
}
