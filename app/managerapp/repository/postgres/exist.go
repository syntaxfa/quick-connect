package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
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
