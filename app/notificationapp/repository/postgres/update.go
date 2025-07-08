package postgres

import (
	"context"
	"encoding/json"

	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryUpdateTemplate = `UPDATE templates
SET name = $1, bodies = $2
WHERE id = $3;`

func (d *DB) UpdateTemplate(ctx context.Context, id types.ID, req service.AddTemplateRequest) error {
	const op = "repository.postgres.update.UpdateTemplate"

	jsonBodies, mErr := json.Marshal(req.Bodies)
	if mErr != nil {
		return richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := d.conn.Conn().Exec(ctx, queryUpdateTemplate, req.Name, jsonBodies, id); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
