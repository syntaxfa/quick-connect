package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryIsExistUserIDFromExternalUserID = `SELECT EXISTS (
	SELECT 1
	FROM external_users
	WHERE external_user_id = $1
);`

func (d *DB) IsExistUserIDFromExternalUserID(ctx context.Context, externalUserID string) (bool, error) {
	const op = "repository.postgres.exist.IsExistUserIDFromExternalUserID"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistUserIDFromExternalUserID, externalUserID).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}

func (d *DB) IsExistTemplateName(_ context.Context, _ string) (bool, error) {
	return false, nil
}
