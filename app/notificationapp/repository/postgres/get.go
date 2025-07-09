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

func (d *DB) FindNotificationByUserID(ctx context.Context, userID types.ID, paginated paginate.RequestBase, isRead *bool) (service.ListNotificationResponse, error) {
	const op = "repository.get.FindNotificationByUserID"

	filters := map[paginate.FilterParameter]paginate.Filter{
		"user_id": {Operation: paginate.FilterOperationEqual, Values: []interface{}{userID}},
	}

	if isRead != nil {
		filters["is_read"] = paginate.Filter{Operation: paginate.FilterOperationEqual, Values: []interface{}{*isRead}}
	}

	fields := []string{"id", "user_id", "type", "title", "body", "data", "is_read", "created_at"}
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
		return service.ListNotificationResponse{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var notifications []service.ListNotificationResult
	for rows.Next() {
		var notification service.ListNotificationResult
		if sErr := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.Title,
			&notification.Body, &notification.Data, &notification.IsRead, &notification.CreatedAt); sErr != nil {
			return service.ListNotificationResponse{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}
		notifications = append(notifications, notification)
	}

	if rErr := rows.Err(); rErr != nil {
		return service.ListNotificationResponse{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return service.ListNotificationResponse{
		Results: notifications,
		Paginate: paginate.ResponseBase{
			CurrentPage: paginated.CurrentPage,
			PageSize:    paginated.PageSize,
		},
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
	var jsonBodies json.RawMessage

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetTemplateByName, name).
		Scan(&template.ID, &template.Name, &jsonBodies, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if uErr := json.Unmarshal(jsonBodies, &template.Bodies); uErr != nil {
		return service.Template{}, richerror.New(op).WithMessage("failed to unmarshalling template bodies").
			WithWrapError(uErr).WithKind(richerror.KindUnexpected)
	}

	return template, nil
}

const queryTemplateByID = `SELECT id, name, bodies, created_at, updated_at
FROM templates WHERE id = $1
LIMIT 1;`

func (d *DB) GetTemplateByID(ctx context.Context, id types.ID) (service.Template, error) {
	const op = "repository.postgres.get.GetTemplateByID"

	var template service.Template
	var jsonBodies json.RawMessage

	if qErr := d.conn.Conn().QueryRow(ctx, queryTemplateByID, id).
		Scan(&template.ID, &template.Name, &jsonBodies, &template.CreatedAt, &template.UpdatedAt); qErr != nil {
		return service.Template{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	if uErr := json.Unmarshal(jsonBodies, &template.Bodies); uErr != nil {
		return service.Template{}, richerror.New(op).WithMessage("failed to unmarshalling template bodies").
			WithWrapError(uErr).WithKind(richerror.KindUnexpected)
	}

	return template, nil
}

func (d *DB) GetUserSetting(_ context.Context, _ types.ID) (service.UserSetting, error) {
	return service.UserSetting{}, nil
}
