package postgres

import (
	"context"
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryCreateNotification = `INSERT INTO notifications (id, user_id, type, title, body, data, channel_deliveries)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, type, title, body, data, is_read, created_at, overall_status, channel_deliveries;`

func (d *DB) Save(ctx context.Context, req service.SendNotificationRequest) (service.Notification, error) {
	const op = "repository.postgres.create.Save"

	var notification service.Notification
	if qErr := d.conn.Conn().QueryRow(ctx, queryCreateNotification, req.ID, req.UserID, req.Type, req.Title, req.Body, req.Data, req.ChannelDeliveries).Scan(
		&notification.ID, &notification.UserID, &notification.Type,
		&notification.Title, &notification.Body, &notification.Data,
		&notification.IsRead, &notification.CreatedAt, &notification.OverallStatus,
		&notification.ChannelDeliveries); qErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't insert into notifications table").WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return notification, nil
}

const queryCreateUserFromExternalUserID = `INSERT INTO external_users (user_id, external_user_id)
VALUES ($1, $2);`

func (d *DB) CreateUserIDFromExternalUserID(ctx context.Context, externalUserID string, userID types.ID) error {
	const op = "repository.postgres.create.CreateUserIDFromExternalUserID"

	_, eErr := d.conn.Conn().Exec(ctx, queryCreateUserFromExternalUserID, userID, externalUserID)
	if eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}

const queryCreateTemplate = `INSERT INTO templates (id, name, bodies)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;`

func (d *DB) CreateTemplate(ctx context.Context, req service.AddTemplateRequest) (service.Template, error) {
	const op = "repository.postgres.create.CreateTemplate"

	jsonBodies, mErr := json.Marshal(req.Bodies)
	if mErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected)
	}

	var template service.Template
	if qErr := d.conn.Conn().QueryRow(ctx, queryCreateTemplate, ulid.Make().String(), req.Name, jsonBodies).
		Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	template.Name = req.Name
	template.Bodies = req.Bodies

	return template, nil
}

const queryCreateUserSetting = `INSERT INTO user_notification_settings (id, user_id, lang, ignore_channels)
VALUES ($1, $2, $3, $4);`

func (d *DB) CreateUserSetting(ctx context.Context, userID types.ID, req service.UpdateUserSettingRequest) (service.UserSetting, error) {
	const op = "repository.postgres.create.CreateUserSetting"

	id := ulid.Make().String()

	jsonChannel, mErr := json.Marshal(req.IgnoreChannels)
	if mErr != nil {
		return service.UserSetting{}, richerror.New(op).WithMessage("can't marshal ignore channels").
			WithWrapError(mErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := d.conn.Conn().Exec(ctx, queryCreateUserSetting, id, userID, req.Lang, jsonChannel); eErr != nil {
		return service.UserSetting{}, richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return service.UserSetting{
		ID:             types.ID(id),
		UserID:         userID,
		Lang:           req.Lang,
		IgnoreChannels: req.IgnoreChannels,
	}, nil
}
