package postgres

import (
	"context"
	"database/sql"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	pagesql "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset/sql"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryGetUserByUsername = `SELECT id, username, hashed_password, fullname, email, phone_number, avatar, last_online_at
FROM users
WHERE username=$1
limit 1;`

const queryGetUserRolesByUserID = `SELECT role
FROM user_roles
WHERE user_id = $1;`

type nullableFields struct {
	Avatar sql.NullString
}

func (d *DB) GetUserByUsername(ctx context.Context, username string) (userservice.User, error) {
	const op = "repository.postgres.GetUserByUsername"

	var user userservice.User
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserByUsername, username).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.Fullname, &user.Email, &user.PhoneNumber, &nullable.Avatar,
		&user.LastOnlineAt); qErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).WithMessage("get user")
	}

	if nullable.Avatar.Valid {
		user.Avatar = nullable.Avatar.String
	}

	roles, grErr := d.GetUserRolesByUserID(ctx, user.ID)
	if grErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(grErr).WithKind(richerror.KindUnexpected)
	}

	user.Roles = roles

	return user, nil
}

const queryGetUserByID = `SELECT id, username, hashed_password, fullname, email, phone_number, avatar, last_online_at
FROM users
WHERE id=$1
limit 1;`

func (d *DB) GetUserByID(ctx context.Context, userID types.ID) (userservice.User, error) {
	const op = "repository.postgres.GetUserByID"

	var user userservice.User
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserByID, userID).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.Fullname, &user.Email, &user.PhoneNumber, &nullable.Avatar,
		&user.LastOnlineAt); qErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).WithMessage("get user")
	}

	if nullable.Avatar.Valid {
		user.Avatar = nullable.Avatar.String
	}

	roles, grErr := d.GetUserRolesByUserID(ctx, user.ID)
	if grErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(grErr).WithKind(richerror.KindUnexpected)
	}

	user.Roles = roles

	return user, nil
}

func (d *DB) GetUserRolesByUserID(ctx context.Context, userID types.ID) ([]types.Role, error) {
	const op = "repository.postgres.GetUserRolesByUserID"

	rows, qrErr := d.conn.Conn().Query(ctx, queryGetUserRolesByUserID, userID)
	if qrErr != nil {
		return nil, richerror.New(op).WithWrapError(qrErr).WithKind(richerror.KindUnexpected).WithMessage("error in Query method for user roles")
	}

	var roles = make([]types.Role, 0)
	for rows.Next() {
		var role types.Role
		if sErr := rows.Scan(&role); sErr != nil {
			return nil, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).WithMessage("error in scan rows user roles")
		}

		roles = append(roles, role)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected).WithMessage("error in rows user roles after scan")
	}

	return roles, nil
}

func (d *DB) GetUserList(ctx context.Context, paginated paginate.RequestBase, username string) ([]userservice.User, paginate.ResponseBase, error) {
	const op = "repository.postgres.get.GetUserList"

	filters := map[paginate.FilterParameter]paginate.Filter{}

	if username != "" {
		filters["username"] = paginate.Filter{Operation: paginate.FilterOperationEqual, Values: []interface{}{username}}
	}

	fields := []string{
		"id", "username", "fullname", "email", "phone_number", "avatar", "last_online_at",
	}
	sortColumn := "id"
	offset := (paginated.CurrentPage - 1) * paginated.PageSize
	limit := paginated.PageSize

	query, countQuery, args := pagesql.WriteQuery(pagesql.Parameters{
		Table:      "users",
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

	var users []userservice.User

	for rows.Next() {
		var user userservice.User
		var nullable nullableFields

		if sErr := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.Email, &user.PhoneNumber,
			&nullable.Avatar, &user.LastOnlineAt); sErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
		}
		
		if nullable.Avatar.Valid {
			user.Avatar = nullable.Avatar.String
		}

		users = append(users, user)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	return users, paginate.ResponseBase{
		CurrentPage: paginated.CurrentPage,
		PageSize:    paginated.PageSize,
	}, nil
}
