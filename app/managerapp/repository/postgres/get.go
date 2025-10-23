package postgres

import (
	"context"
	"database/sql"
	"github.com/syntaxfa/quick-connect/types"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
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

	rows, qrErr := d.conn.Conn().Query(ctx, queryGetUserRolesByUserID, user.ID)
	if qrErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qrErr).WithKind(richerror.KindUnexpected).WithMessage("error in Query method for user roles")
	}
	for rows.Next() {
		var role types.Role
		if sErr := rows.Scan(&role); sErr != nil {
			return userservice.User{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).WithMessage("error in scan rows user roles")
		}

		user.Roles = append(user.Roles, role)
	}

	if rErr := rows.Err(); rErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected).WithMessage("error in rows user roles after scan")
	}

	return user, nil
}
