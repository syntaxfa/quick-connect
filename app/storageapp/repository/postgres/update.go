package postgres

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryDeleteByID = `UPDATE files
SET deleted_at = $1
WHERE id = $2;`

func (d *DB) DeleteByID(ctx context.Context, fileID types.ID) error {
	const op = "repository.postgres.update.DeleteByID"

	cmdTag, exErr := d.conn.Conn().Exec(ctx, queryDeleteByID, time.Now(), fileID)
	if exErr != nil {
		return richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected)
	}

	if cmdTag.RowsAffected() == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage("file not found for delete")
	}

	return nil
}

const queryConfirmFile = `UPDATE files
SET is_confirmed = true
WHERE id = $1;`

func (d *DB) ConfirmFile(ctx context.Context, fileID types.ID) error {
	const op = "repository.postgres.update.ConfirmFile"

	cmdTag, exErr := d.conn.Conn().Exec(ctx, queryConfirmFile, fileID)
	if exErr != nil {
		return richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected)
	}

	if cmdTag.RowsAffected() == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage("file not found for delete")
	}

	return nil
}
