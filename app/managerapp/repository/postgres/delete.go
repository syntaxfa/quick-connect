package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryDeleteUser = `DELETE FROM users
WHERE id = $1;`

func (d *DB) DeleteUser(ctx context.Context, userID types.ID) error {
	const op = "repository.postgres.delete.DeleteUser"

	if _, eErr := d.conn.Conn().Exec(ctx, queryDeleteUser, userID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}

const queryDeleteUserRole = `DELETE FROM user_roles
WHERE user_id = $1;`

func (d *DB) DeleteUserRole(ctx context.Context, userID types.ID) error {
	const op = "repository.postgres.delete.DeleteUserRole"

	if _, eErr := d.conn.Conn().Exec(ctx, queryDeleteUserRole, userID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
