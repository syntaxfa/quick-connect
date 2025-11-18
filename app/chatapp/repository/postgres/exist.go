package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryIsUserHaveActiveConversation = `SELECT EXISTS (
	SELECT 1
	FROM conversations
	WHERE client_user_id = $1 AND status IN ('new', 'open', 'bot_handling')
);`

func (d *DB) IsUserHaveActiveConversation(ctx context.Context, userID types.ID) (bool, error) {
	const op = "repository.postgres.exist.IsUserHaveActiveConversation"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryIsUserHaveActiveConversation, userID).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}

		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}

const queryCheckUserInConversation = `
SELECT EXISTS (
    SELECT 1
    FROM conversations
    WHERE id = $1 AND (client_user_id = $2 OR assigned_support_id = $2)
);`

func (d *DB) CheckUserInConversation(ctx context.Context, userID, conversationID types.ID) (bool, error) {
	const op = "repository.postgres.exist.CheckUserInConversation"

	var exists bool
	if qErr := d.conn.Conn().QueryRow(ctx, queryCheckUserInConversation, conversationID, userID).Scan(&exists); qErr != nil {
		if errors.Is(qErr, pgx.ErrNoRows) {
			return false, nil
		}
		return false, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected)
	}

	return exists, nil
}
