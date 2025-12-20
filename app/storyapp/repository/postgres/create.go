package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/storyapp/service"
)

func (d *DB) SaveStory(_ context.Context, _ service.AddStoryRequest) (service.Story, error) {
	return service.Story{}, nil
}
