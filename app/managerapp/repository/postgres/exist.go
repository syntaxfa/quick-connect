package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryIsExistUserByUsername = `SELECT EXISTS (
	SELECT 1
	FROM users
	WHERE username = $1
);`

func (d *DB) IsExistUserByUsername(ctx context.Context, username string) (bool, error) {
	const op = "repository.postgres.IsExistUserByUsername"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistUserByUsername, username).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}

const queryIsExistUserByID = `SELECT EXISTS (
	SELECT 1
	FROM users
	WHERE id = $1
);`

func (d *DB) IsExistUserByID(ctx context.Context, userID types.ID) (bool, error) {
	const op = "repository.postgres.IsExistUserByID"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistUserByID, userID).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}

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
