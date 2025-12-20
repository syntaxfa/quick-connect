package service

import "context"

func (s Service) AddStory(_ context.Context, req AddStoryRequest) (AddStoryResponse, error) {
	const op = "service.add_story.AddStory"

	_ = op

	if vErr := s.vld.ValidateAddStoryRequest(req); vErr != nil {
		return AddStoryResponse{}, vErr
	}

	return AddStoryResponse{}, nil
}
