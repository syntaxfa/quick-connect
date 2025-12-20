package service

import "context"

type Repository interface {
	SaveStory(ctx context.Context, req AddStoryRequest) (Story, error)
}

type Service struct {
	repo Repository
	vld  Validate
}

func New(repo Repository, vld Validate) Service {
	return Service{
		repo: repo,
		vld:  vld,
	}
}
