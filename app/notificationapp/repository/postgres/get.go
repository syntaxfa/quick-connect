package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (d *DB) FindNotificationByUserID(_ context.Context, _ types.ID) ([]service.Notification, error) {
	// TODO: add paginated
	return nil, nil
}

const queryGetUserIDFromExternalUserID = `SELECT user_id FROM external_users
WHERE external_user_id = $1
LIMIT 1;`

func (d *DB) GetUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error) {
	const op = "repository.postgres.get.GetUserIDFromExternalUserID"

	var userID string
	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserIDFromExternalUserID, externalUserID).Scan(&userID); qErr != nil {
		return "", richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return types.ID(userID), nil
}
