package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
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
	return d.getUserBy(ctx, op, queryGetUserByUsername, username)
}

const queryGetUserByID = `SELECT id, username, hashed_password, fullname, email, phone_number, avatar, last_online_at
FROM users
WHERE id=$1
limit 1;`

func (d *DB) GetUserByID(ctx context.Context, userID types.ID) (userservice.User, error) {
	const op = "repository.postgres.GetUserByID"
	return d.getUserBy(ctx, op, queryGetUserByID, userID)
}

// getUserBy is a private helper function that encapsulates the duplicated user fetching logic.
func (d *DB) getUserBy(ctx context.Context, op string, query string, arg interface{}) (userservice.User, error) {
	var user userservice.User
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, query, arg).Scan(
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
		return nil, richerror.New(op).WithWrapError(qrErr).WithKind(richerror.KindUnexpected).
			WithMessage("error in Query method for user roles")
	}

	var roles = make([]types.Role, 0)
	for rows.Next() {
		var role types.Role
		if sErr := rows.Scan(&role); sErr != nil {
			return nil, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
				WithMessage("error in scan rows user roles")
		}

		roles = append(roles, role)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected).
			WithMessage("error in rows user roles after scan")
	}

	return roles, nil
}

func (d *DB) buildUserListFilters(username string, roles []types.Role) (string, string, []interface{}) {
	args := make([]interface{}, 0)
	whereClauses := make([]string, 0)
	joinClause := ""
	argCount := 1

	if username != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.username ILIKE $%d", argCount))
		args = append(args, "%"+username+"%")
		argCount++
	}

	if len(roles) > 0 {
		joinClause = "JOIN user_roles ur ON u.id = ur.user_id"
		whereClauses = append(whereClauses, fmt.Sprintf("ur.role = ANY($%d)", argCount))
		args = append(args, pq.Array(roles))
	}

	whereQuery := ""
	if len(whereClauses) > 0 {
		whereQuery = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	return joinClause, whereQuery, args
}

func (d *DB) GetUserList(ctx context.Context, paginated paginate.RequestBase,
	username string, roles []types.Role) ([]userservice.User, paginate.ResponseBase, error) {
	const op = "repository.postgres.get.GetUserList"

	offset := (paginated.CurrentPage - 1) * paginated.PageSize
	limit := paginated.PageSize
	sortColumn := "id"
	sortDirection := "ASC"
	if paginated.Descending {
		sortDirection = "DESC"
	}

	joinClause, whereQuery, args := d.buildUserListFilters(username, roles)

	const (
		limitArgDelta  = 1
		offsetArgDelta = 2
	)

	argCount := len(args)
	query := fmt.Sprintf(`
		SELECT DISTINCT u.id, u.username, u.fullname, u.email, u.phone_number, u.avatar, u.last_online_at
		FROM users u
		%s
		%s
		ORDER BY u.%s %s
		LIMIT $%d OFFSET $%d`,
		joinClause, whereQuery, sortColumn, sortDirection, argCount+limitArgDelta, argCount+offsetArgDelta)

	mainQueryArgs := make([]interface{}, 0, len(args)+offsetArgDelta)
	mainQueryArgs = append(mainQueryArgs, args...)
	mainQueryArgs = append(mainQueryArgs, limit, offset)

	rows, qErr := d.conn.Conn().Query(ctx, query, mainQueryArgs...)
	if qErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).
			WithMessage("query error")
	}
	defer rows.Close()

	var users []userservice.User
	for rows.Next() {
		var user userservice.User
		var nullable nullableFields

		if sErr := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.Email, &user.PhoneNumber, &nullable.Avatar,
			&user.LastOnlineAt); sErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
				WithMessage("scan error")
		}

		if nullable.Avatar.Valid {
			user.Avatar = nullable.Avatar.String
		}

		users = append(users, user)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		%s
		%s`,
		joinClause, whereQuery)

	var totalCount uint64
	if sErr := d.conn.Conn().QueryRow(ctx, countQuery, args...).Scan(&totalCount); sErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
			WithMessage("count query error")
	}

	return users, paginate.ResponseBase{
		CurrentPage:  paginated.CurrentPage,
		PageSize:     paginated.PageSize,
		TotalNumbers: totalCount,
		TotalPage:    (totalCount + paginated.PageSize - 1) / paginated.PageSize,
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
