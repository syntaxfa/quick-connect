package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const queryCreateUser = `INSERT INTO users (id, username, hashed_password, fullname)
VALUES ($1, $2, $3, $4)
RETURNING id, username, fullname, last_online_at;`

const queryCreateUserRole = `INSERT INTO user_roles (user_id, role)
VALUES ($1, $2);`

func (d *DB) CreateUser(ctx context.Context, req userservice.UserCreateRequest) (userservice.User, error) {
	const op = "repository.postgres.create.CreateUser"

	tx, tErr := d.conn.Conn().Begin(ctx)
	if tErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(tErr).WithKind(richerror.KindUnexpected).
			WithMessage("error while begin a transaction")
	}

	var user userservice.User

	if sErr := tx.QueryRow(ctx, queryCreateUser, req.ID, req.Username, req.Password, req.Fullname).
		Scan(&user.ID, &user.Username, &user.Fullname, &user.LastOnlineAt); sErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
	}

	for _, role := range req.Roles {
		if _, srErr := tx.Exec(ctx, queryCreateUserRole, user.ID, role); srErr != nil {
			return userservice.User{}, richerror.New(op).WithWrapError(srErr).WithKind(richerror.KindUnexpected)
		}

		user.Roles = append(user.Roles, role)
	}

	if cErr := tx.Commit(ctx); cErr != nil {
		return userservice.User{}, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected).
			WithMessage("error while transaction commit")
	}

	return user, nil
}
