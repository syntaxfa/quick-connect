package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/types"
)

func (d *DB) FindByUserID(_ context.Context, _ types.ID) ([]service.Notification, error) {
	return nil, nil
}

func (d *DB) GetUserIDFromExternalUserID(_ context.Context, _ string) (types.ID, error) {
	return "", nil
}
