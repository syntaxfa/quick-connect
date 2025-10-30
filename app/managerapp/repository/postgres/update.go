package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryUpdateUser = `UPDATE users
SET username = $1, fullname = $2, email = $3, phone_number = $4
WHERE id = $5;`

func (d *DB) UpdateUser(ctx context.Context, userID types.ID, req userservice.UserUpdateFromSuperuserRequest) error {
	const op = "repository.postgres.update.UpdateUser"

	tx, bErr := d.conn.Conn().Begin(ctx)
	if bErr != nil {
		return richerror.New(op).WithWrapError(bErr).WithKind(richerror.KindUnexpected)
	}

	if _, eErr := d.conn.Conn().Exec(ctx, queryUpdateUser, req.Username, req.Fullname, req.Email, &req.PhoneNumber, userID); eErr != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			return richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
		}
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	if delErr := d.DeleteUserRole(ctx, userID); delErr != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			return richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
		}
		return richerror.New(op).WithWrapError(delErr).WithKind(richerror.KindUnexpected)
	}

	for _, role := range req.Roles {
		if _, srErr := tx.Exec(ctx, queryCreateUserRole, userID, role); srErr != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				return richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
			}
			return richerror.New(op).WithWrapError(srErr).WithKind(richerror.KindUnexpected)
		}
	}

	if cErr := tx.Commit(ctx); cErr != nil {
		return richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
