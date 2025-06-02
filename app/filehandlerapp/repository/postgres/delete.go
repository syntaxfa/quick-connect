package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const deleteQuery = "UPDATE files SET is_deleted = true WHERE id = $1;"

func (d *DB) Delete(ctx context.Context, id types.ULID) error {
	const op = "repository.postgres.Delete"

	if _, qErr := d.conn.Conn().Exec(ctx, createQuery, id); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return nil
		}

		return richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
