package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
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

const queryIsExistTemplateByName = `SELECT EXISTS (
	SELECT 1
	FROM templates
	WHERE name = $1
);`

func (d *DB) IsExistTemplateByName(ctx context.Context, name string) (bool, error) {
	const op = "repository.postgres.exist.IsExistTemplateByName"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistTemplateByName, name).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}

const queryIsExistTemplateByID = `SELECT EXISTS (
	SELECT 1
	FROM templates
	WHERE id = $1
);`

func (d *DB) IsExistTemplateByID(ctx context.Context, id types.ID) (bool, error) {
	const op = "repository.postgres.exist.IsExistTemplateByID"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistTemplateByID, id).Scan(exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}
