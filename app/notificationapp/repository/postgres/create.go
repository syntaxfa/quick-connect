package postgres

import (
	"context"
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryCreateNotification = `INSERT INTO notifications (id, user_id, type, data, template_name, dynamic_body_data, dynamic_title_data, is_in_app, overall_status, channel_deliveries)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, user_id, type, data, template_name, dynamic_body_data, dynamic_title_data, is_read, is_in_app, created_at, overall_status, channel_deliveries;`

func (d *DB) Save(ctx context.Context, req service.SendNotificationRequest) (service.Notification, error) {
	const op = "repository.postgres.create.Save"

	jsonData, mdErr := json.Marshal(req.Data)
	if mdErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't marshal notification data").
			WithWrapError(mdErr).WithKind(richerror.KindUnexpected)
	}

	jsonBodyData, mbErr := json.Marshal(req.DynamicBodyData)
	if mbErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't marshal notification dynamic body data").
			WithWrapError(mbErr).WithKind(richerror.KindUnexpected)
	}

	jsonTitleData, mtErr := json.Marshal(req.DynamicTitleData)
	if mtErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't marshal notification dynamic title data").
			WithWrapError(mtErr).WithKind(richerror.KindUnexpected)
	}

	var notification service.Notification
	var jsonChannelDeliveries json.RawMessage
	if qErr := d.conn.Conn().QueryRow(ctx, queryCreateNotification, req.ID, req.UserID, req.Type, jsonData, req.TemplateName,
		jsonBodyData, jsonTitleData, req.IsInApp, req.Status, req.ChannelDeliveries).Scan(
		&notification.ID, &notification.UserID, &notification.Type, &jsonData, &notification.TemplateName, &jsonBodyData,
		&jsonTitleData, &notification.IsRead, &notification.IsInApp, &notification.CreatedAt, &notification.OverallStatus,
		&jsonChannelDeliveries); qErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't insert into notifications table").
			WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if udErr := json.Unmarshal(jsonData, &notification.Data); udErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't unmarshall notification data").
			WithWrapError(udErr).WithKind(richerror.KindUnexpected)
	}

	if ubErr := json.Unmarshal(jsonBodyData, &notification.DynamicBodyData); ubErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't unmarshall notification dynamic body data").
			WithWrapError(ubErr).WithKind(richerror.KindUnexpected)
	}

	if utErr := json.Unmarshal(jsonTitleData, &notification.DynamicTitleData); utErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't unmarshal notification title body data").
			WithWrapError(utErr).WithKind(richerror.KindUnexpected)
	}

	if ucErr := json.Unmarshal(jsonChannelDeliveries, &notification.ChannelDeliveries); ucErr != nil {
		return service.Notification{}, richerror.New(op).WithMessage("can't unmarshal notification channel deliveries").
			WithWrapError(ucErr).WithKind(richerror.KindUnexpected)
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

const queryCreateTemplate = `INSERT INTO templates (id, name, contents)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;`

func (d *DB) CreateTemplate(ctx context.Context, req service.AddTemplateRequest) (service.Template, error) {
	const op = "repository.postgres.create.CreateTemplate"

	jsonContents, mErr := json.Marshal(req.Contents)
	if mErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected)
	}

	var template service.Template
	if qErr := d.conn.Conn().QueryRow(ctx, queryCreateTemplate, req.ID, req.Name, jsonContents).
		Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	template.Name = req.Name
	template.Contents = req.Contents

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
