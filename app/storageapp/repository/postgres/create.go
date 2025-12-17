package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const querySave = `INSERT INTO files (id, uploader_id, name, key, mime_type, size, driver, bucket, is_public)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

func (d *DB) Save(ctx context.Context, file service.File) error {
	const op = "repository.postgres.create.Save"

	var nullable nullableFields
	if file.Bucket != "" {
		nullable.Bucket.String = file.Bucket
		nullable.Bucket.Valid = true
	}

	if _, exErr := d.conn.Conn().Exec(ctx, querySave, file.ID, file.UploaderID, file.Name, file.Key, file.MimeType,
		file.Size, file.Driver, nullable.Bucket, file.IsPublic); exErr != nil {
		return richerror.New(op).WithWrapError(exErr).WithMessage("can't insert file").WithKind(richerror.KindUnexpected)
	}

	return nil
}
