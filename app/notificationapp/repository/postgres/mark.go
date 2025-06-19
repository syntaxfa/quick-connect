package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryMarkAllAsRead = `UPDATE notifications
SET is_read = true
WHERE user_id = $1;`

func (d *DB) MarkAllAsReadByUserID(ctx context.Context, userID types.ID) error {
	const op = "repository.mark.MarkAllAsReadByUserID"

	_, eErr := d.conn.Conn().Exec(ctx, queryMarkAllAsRead, userID)
	if eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}

const queryMarkAsRead = `UPDATE notifications
SET is_read = true
WHERE id = $1 AND user_id = $2;`

func (d *DB) MarkAsRead(ctx context.Context, notificationID, userID types.ID) error {
	const op = "repository.mark.MarkAsRead"

	_, eErr := d.conn.Conn().Exec(ctx, queryMarkAsRead, notificationID, userID)
	if eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
