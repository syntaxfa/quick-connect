package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryIsExistByID = `SELECT EXISTS (
	SELECT 1
	FROM files
	WHERE id = $1
);`

func (d *DB) IsExistByID(ctx context.Context, fileID types.ID) (bool, error) {
	const op = "repository.postgres.exist.IsExistByID"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsExistByID, fileID).Scan(&exists); qErr != nil {
		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}
