package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

//nolint:gosec // G101: This is a SQL query template, not a hardcoded credential
const queryPasswordIsCorrect = `SELECT EXISTS (
	SELECT 1
	FROM users
	WHERE id = $1 AND hashed_password = $2
);`

func (d *DB) PasswordIsCorrect(ctx context.Context, userID types.ID, hashedPassword string) (bool, error) {
	const op = "repository.postgres.password.PasswordIsCorrect"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryPasswordIsCorrect, userID, hashedPassword).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return true, nil
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
