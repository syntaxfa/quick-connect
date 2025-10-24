package postgres

import (
	"context"
	"database/sql"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryGetUserByUsername = `SELECT id, username, hashed_password, fullname, avatar, last_online_at
FROM users
WHERE username=$1
limit 1;`

const queryGetUserRolesByUserID = `SELECT role
FROM user_roles
WHERE user_id = $1;`

type nullableFields struct {
	Fullname sql.NullString
	Avatar   sql.NullString
}

func (d *DB) GetUserByUsername(ctx context.Context, username string) (userservice.User, error) {
	const op = "repository.postgres.GetUserByUsername"

	var user userservice.User
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserByUsername, username).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &nullable.Fullname, &nullable.Avatar,
		&user.LastOnlineAt); qErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).WithMessage("get user")
	}

	if nullable.Fullname.Valid {
		user.Fullname = nullable.Fullname.String
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

const queryGetUserByID = `SELECT id, username, hashed_password, fullname, avatar, last_online_at
FROM users
WHERE id=$1
limit 1;`

func (d *DB) GetUserByID(ctx context.Context, userID types.ID) (userservice.User, error) {
	const op = "repository.postgres.GetUserByID"

	var user userservice.User
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserByID, userID).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &nullable.Fullname, &nullable.Avatar,
		&user.LastOnlineAt); qErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).WithMessage("get user")
	}

	if nullable.Fullname.Valid {
		user.Fullname = nullable.Fullname.String
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
