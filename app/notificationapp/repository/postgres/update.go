package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
)

func (d *DB) UpdateTemplate(_ context.Context, _ service.AddTemplateRequest) (service.AddTemplateResponse, error) {
	return service.AddTemplateResponse{}, nil
}
