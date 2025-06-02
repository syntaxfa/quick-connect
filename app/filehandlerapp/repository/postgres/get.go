package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/app/filehandlerapp/service/file"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const getQuery = `
SELECT id, type, type_id, extension, storage_type, size, content_type, created_at FROM files 
WHERE id = $1 AND is_deleted = false
`

func (d *DB) Get(ctx context.Context, id types.ULID) (file.File, error) {
	const op = "repository.postgres.Create"

	var f file.File
	qErr := d.conn.Conn().QueryRow(ctx, createQuery, id).Scan(
		&id, &f.Type, &f.TypeID, &f.Extension, &f.StorageType, &f.Size, &f.ContentType, &f.CreatedAt,
	)
	if qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return file.File{}, nil
		}

		return file.File{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return f, nil
}
