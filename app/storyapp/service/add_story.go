package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) AddStory(ctx context.Context, req AddStoryRequest, creatorID types.ID) (AddStoryResponse, error) {
	const op = "service.add_story.AddStory"

	if vErr := s.vld.ValidateAddStoryRequest(req); vErr != nil {
		return AddStoryResponse{}, vErr
	}

	ctxWithValue, tErr := s.tokenManager.SetTokenInContext(ctx)
	if tErr != nil {
		return AddStoryResponse{}, errlog.ErrContext(ctxWithValue, richerror.New(op).WithWrapError(tErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	filePb, getFErr := s.storageSvc.GetFileInfo(ctxWithValue, &storagepb.GetFileInfoRequest{FileId: string(req.MediaFileID)})
	if getFErr != nil {
		return AddStoryResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(getFErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if filePb.GetIsConfirmed() {
		return AddStoryResponse{}, richerror.New(op).WithMessage(servermsg.MsgMediaAlreadyUse)
	}

	if !filePb.GetIsPublic() {
		return AddStoryResponse{}, richerror.New(op).WithMessage(servermsg.MsgStoryMediaRequirePublic).WithKind(richerror.KindBadRequest)
	}

	// confirmed file.

	req.CreatorID = creatorID
	story, saveErr := s.repo.SaveStory(ctx, req)
	if saveErr != nil {
		return AddStoryResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(saveErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return AddStoryResponse{
		Story: story,
	}, nil
}
