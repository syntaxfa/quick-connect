package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/app/filehandlerapp/service/file"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const createQuery = `
INSERT INTO files (type, type_id, extension, storage_type, size, content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`

func (d *DB) Create(ctx context.Context, file file.File) (types.ULID, error) {
	const op = "repository.postgres.Create"

	var id types.ULID
	qErr := d.conn.Conn().QueryRow(
		ctx, createQuery, file.Type, file.TypeID, file.Extension, file.StorageType, file.Size, file.ContentType,
	).Scan(&id)
	if qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return "", nil
		}

		return "", richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return id, nil
}
