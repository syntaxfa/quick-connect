package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

//nolint:gosec // G101: This is a SQL query template, not a hardcoded credential
const queryGetUserHashedPassword = `SELECT hashed_password
FROM users
WHERE id = $1
limit 1;`

func (d *DB) GetUserHashedPassword(ctx context.Context, userID types.ID) (string, error) {
	const op = "repository.postgres.password.GetUserHashedPassword"

	var hashedPass string
	if qErr := d.conn.Conn().QueryRow(ctx, queryGetUserHashedPassword, userID).Scan(&hashedPass); qErr != nil {
		return "", richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return hashedPass, nil
}

//nolint:gosec // G101: This is a SQL query template, not a hardcoded credential
const queryChangePassword = `UPDATE users
SET hashed_password = $1
WHERE id = $2;`

func (d *DB) ChangePassword(ctx context.Context, userID types.ID, hashedPassword string) error {
	const op = "repository.postgres.password.ChangePassword"

	if _, eErr := d.conn.Conn().Exec(ctx, queryChangePassword, hashedPassword, userID); eErr != nil {
		return richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected)
	}

	return nil
}
