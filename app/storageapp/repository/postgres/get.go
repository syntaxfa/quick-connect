package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryGetByID = `SELECT id, uploader_id, name, key, mime_type, size, driver,
bucket, is_public, is_confirmed, created_at, updated_at, deleted_at
FROM files
WHERE id = $1
LIMIT 1;`

func (d *DB) GetByID(ctx context.Context, fileID types.ID) (service.File, error) {
	const op = "repository.postgres.get.GetByID"

	var file service.File
	var nullable nullableFields

	if sErr := d.conn.Conn().QueryRow(ctx, queryGetByID, fileID).Scan(&file.ID, &file.UploaderID, &file.Name,
		&file.Key, &file.MimeType, &file.Size, &file.Driver, &nullable.Bucket, &file.IsPublic, &file.IsConfirmed,
		&file.CreatedAt, &file.UpdatedAt, &nullable.DeletedAt); sErr != nil {
		return service.File{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
	}

	if nullable.Bucket.Valid {
		file.Bucket = nullable.Bucket.String
	}
	if nullable.DeletedAt.Valid {
		file.DeletedAt = &nullable.DeletedAt.Time
	}

	return file, nil
}
