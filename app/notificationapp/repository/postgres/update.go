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

const queryUpdateUserSetting = `UPDATE user_notification_settings
SET lang = $1, ignore_channels = $2
WHERE user_id = $3;`

func (d *DB) UpdateUserSetting(ctx context.Context, userID types.ID, req service.UpdateUserSettingRequest) error {
	const op = "repository.postgres.update.UpdateUserSetting"

	jsonChannels, mErr := json.Marshal(req.IgnoreChannels)
	if mErr != nil {
		return richerror.New(op).WithMessage("can't marshal ignore channels").WithWrapError(mErr).
			WithKind(richerror.KindUnexpected)
	}

	if _, eErr := d.conn.Conn().Exec(ctx, queryUpdateUserSetting, req.Lang, jsonChannels, userID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
