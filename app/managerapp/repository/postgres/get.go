package postgres

import (
	"context"
	"database/sql"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryGetUserByUsername = `SELECT id, username, hashed_password, fullname, avatar, role, last_online_at
FROM users
WHERE username=$1
limit 1;`

type nullableFields struct {
	Fullname sql.NullString
	Avatar   sql.NullString
}

func (d *DB) GetUserByUsername(ctx context.Context, username string) (userservice.User, error) {
	const op = "repository.postgres.GetUserByUsername"

	var user userservice.User
	var roleString string
	var nullable nullableFields

	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserByUsername, username).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &nullable.Fullname, &nullable.Avatar,
		&roleString, &user.LastOnlineAt); qErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	role, rErr := userservice.RoleStringToInt(roleString)
	if rErr != nil {
		return userservice.User{}, rErr
	}

	user.Role = role

	if nullable.Fullname.Valid {
		user.Fullname = nullable.Fullname.String
	}
	if nullable.Avatar.Valid {
		user.Avatar = nullable.Avatar.String
	}

	return user, nil
}
