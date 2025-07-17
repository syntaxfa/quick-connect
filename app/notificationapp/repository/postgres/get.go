package postgres

import (
	"context"
	"encoding/json"

	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	pagesql "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset/sql"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (d *DB) FindNotificationByUserID(ctx context.Context, userID types.ID, paginated paginate.RequestBase, isRead, isInApp *bool) ([]service.Notification, paginate.ResponseBase, error) {
	const op = "repository.get.FindNotificationByUserID"

	filters := map[paginate.FilterParameter]paginate.Filter{
		"user_id": {Operation: paginate.FilterOperationEqual, Values: []interface{}{userID}},
	}

	if isRead != nil {
		filters["is_read"] = paginate.Filter{Operation: paginate.FilterOperationEqual, Values: []interface{}{*isRead}}
	}

	if isInApp != nil {
		filters["is_in_app"] = paginate.Filter{Operation: paginate.FilterOperationEqual, Values: []interface{}{*isInApp}}
	}

	fields := []string{
		"id", "user_id", "type", "data", "template_name", "dynamic_body_data", "dynamic_title_data",
		"is_read", "created_at", "overall_status", "channel_deliveries",
	}
	sortColumn := "created_at"
	offset := (paginated.CurrentPage - 1) * paginated.PageSize
	limit := paginated.PageSize

	query, countQuery, args := pagesql.WriteQuery(pagesql.Parameters{
		Table:      "notifications",
		Fields:     fields,
		Filters:    filters,
		SortColumn: sortColumn,
		Descending: paginated.Descending,
		Limit:      limit,
		Offset:     offset,
	})

	// TODO: complete this
	_ = countQuery

	rows, qErr := d.conn.Conn().Query(ctx, query, args...)
	if qErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var notifications []service.Notification
	var jsonData json.RawMessage
	var jsonBodyData json.RawMessage
	var jsonTitleData json.RawMessage
	var jsonChannelDelivery json.RawMessage
	for rows.Next() {
		var notification service.Notification
		if sErr := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &jsonData, &notification.TemplateName,
			&jsonBodyData, &jsonTitleData, &notification.IsRead, &notification.CreatedAt, &notification.OverallStatus, &jsonChannelDelivery); sErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}

		if uErr := json.Unmarshal(jsonData, &notification.Data); uErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected).
				WithMessage("can't unmarshal notification data")
		}

		if uErr := json.Unmarshal(jsonBodyData, &notification.DynamicBodyData); uErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected).
				WithMessage("can't unmarshal notification dynamic body data")
		}

		if uErr := json.Unmarshal(jsonTitleData, &notification.DynamicTitleData); uErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected).
				WithMessage("can't unmarshal notification dynamic title data")
		}

		if uErr := json.Unmarshal(jsonChannelDelivery, &notification.ChannelDeliveries); uErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected).
				WithMessage("can't unmarshal notification channel deliveries")
		}

		notifications = append(notifications, notification)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return notifications, paginate.ResponseBase{
		CurrentPage: paginated.CurrentPage,
		PageSize:    paginated.PageSize,
	}, nil
}

const queryGetUserIDFromExternalUserID = `SELECT user_id FROM external_users
WHERE external_user_id = $1
LIMIT 1;`

func (d *DB) GetUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error) {
	const op = "repository.postgres.get.GetUserIDFromExternalUserID"

	var userID string
	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserIDFromExternalUserID, externalUserID).Scan(&userID); qErr != nil {
		return "", richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return types.ID(userID), nil
}

const queryGetTemplateByName = `SELECT id, name, bodies, created_at, updated_at
FROM templates WHERE name = $1
LIMIT 1;`

func (d *DB) GetTemplateByName(ctx context.Context, name string) (service.Template, error) {
	const op = "repository.postgres.get.GetTemplateByName"

	var template service.Template
	var jsonContents json.RawMessage

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetTemplateByName, name).
		Scan(&template.ID, &template.Name, &jsonContents, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if uErr := json.Unmarshal(jsonContents, &template.Contents); uErr != nil {
		return service.Template{}, richerror.New(op).WithMessage("failed to unmarshalling template bodies").
			WithWrapError(uErr).WithKind(richerror.KindUnexpected)
	}

	return template, nil
}

const queryTemplateByID = `SELECT id, name, contents, created_at, updated_at
FROM templates WHERE id = $1
LIMIT 1;`

func (d *DB) GetTemplateByID(ctx context.Context, id types.ID) (service.Template, error) {
	const op = "repository.postgres.get.GetTemplateByID"

	var template service.Template
	var jsonContents json.RawMessage

	if qErr := d.conn.Conn().QueryRow(ctx, queryTemplateByID, id).
		Scan(&template.ID, &template.Name, &jsonContents, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if uErr := json.Unmarshal(jsonContents, &template.Contents); uErr != nil {
		return service.Template{}, richerror.New(op).WithMessage("failed to unmarshalling template bodies").
			WithWrapError(uErr).WithKind(richerror.KindUnexpected)
	}

	return template, nil
}

const queryGetTemplatesByNames = `SELECT id, name, contents, created_at
FROM templates WHERE name = ANY($1)`

func (d *DB) GetTemplatesByNames(ctx context.Context, names ...string) ([]service.Template, error) {
	const op = "repository,postgres.get.GetTemplatesByNames"

	var templates []service.Template

	rows, qErr := d.conn.Conn().Query(ctx, queryGetTemplatesByNames, names)
	if qErr != nil {
		return nil, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		var template service.Template
		var jsonContents json.RawMessage
		if sErr := rows.Scan(&template.ID, &template.Name, &jsonContents, &template.CreatedAt); sErr != nil {
			return nil, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}

		if uErr := json.Unmarshal(jsonContents, &template.Contents); uErr != nil {
			return nil, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected)
		}

		templates = append(templates, template)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return templates, nil
}

func (d *DB) GetTemplates(ctx context.Context, req service.ListTemplateRequest) (service.ListTemplateResponse, error) {
	const op = "repository.get.GetTemplates"

	filters := make(map[paginate.FilterParameter]paginate.Filter)
	if req.Name != "" {
		filters["name"] = paginate.Filter{Operation: paginate.FilterOperationEqual, Values: []interface{}{req.Name}}
	}

	fields := []string{"id", "name", "created_at", "updated_at"}
	sortColumn := "created_at"
	offset := (req.Paginated.CurrentPage - 1) * req.Paginated.PageSize
	limit := req.Paginated.PageSize

	query, countQuery, args := pagesql.WriteQuery(pagesql.Parameters{
		Table:      "templates",
		Fields:     fields,
		Filters:    filters,
		SortColumn: sortColumn,
		Descending: req.Paginated.Descending,
		Limit:      limit,
		Offset:     offset,
	})

	// TODO: complete this
	_ = countQuery

	rows, qErr := d.conn.Conn().Query(ctx, query, args...)
	if qErr != nil {
		return service.ListTemplateResponse{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var templates []service.ListTemplateResult
	for rows.Next() {
		var template service.ListTemplateResult
		if sErr := rows.Scan(&template.ID, &template.Name, &template.CreatedAt, &template.UpdatedAt); sErr != nil {
			return service.ListTemplateResponse{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}
		templates = append(templates, template)
	}

	if rErr := rows.Err(); rErr != nil {
		return service.ListTemplateResponse{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return service.ListTemplateResponse{
		Results: templates,
		Paginate: paginate.ResponseBase{
			CurrentPage: req.Paginated.CurrentPage,
			PageSize:    req.Paginated.PageSize,
		},
	}, nil
}

const queryGetUserSetting = `SELECT id, user_id, lang, ignore_channels
FROM user_notification_settings
WHERE user_id = $1`

func (d *DB) GetUserSetting(ctx context.Context, userID types.ID) (service.UserSetting, error) {
	const op = "repository.postgres.get.GetUserSetting"

	var setting service.UserSetting
	var jsonChannel json.RawMessage
	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserSetting, userID).
		Scan(&setting.ID, &setting.UserID, &setting.Lang, &jsonChannel); qErr != nil {
		return service.UserSetting{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if uErr := json.Unmarshal(jsonChannel, &setting.IgnoreChannels); uErr != nil {
		return service.UserSetting{}, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected)
	}

	return setting, nil
}
