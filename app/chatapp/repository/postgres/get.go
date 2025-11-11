package postgres

import (
	"context"
	"database/sql"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryGetUserActiveConversation = `SELECT id, client_user_id, assigned_support_id, status, last_message_snippet,
last_message_sender_id, created_at, updated_at, closed_at
FROM conversations
WHERE client_user_id = $1 AND status IN ('new', 'open', 'bot_handling')
limit 1;`

type nullableFields struct {
	AssignedSupportID   sql.NullString
	LastMessageSnippet  sql.NullString
	LastMessageSenderID sql.NullString
}

func (d *DB) GetUserActiveConversation(ctx context.Context, userID types.ID) (service.Conversation, error) {
	const op = "repository.postgres.get.GetUserActiveConversation"

	return d.GetConversationBy(ctx, op, queryGetUserActiveConversation, userID)
}

func (d *DB) GetConversationBy(ctx context.Context, op string, query string, arg interface{}) (service.Conversation, error) {
	var conversation service.Conversation
	var nullable nullableFields

	if sErr := d.conn.Conn().QueryRow(ctx, query, arg).Scan(&conversation.ID, &conversation.ClientUserID, &nullable.AssignedSupportID,
		&conversation.Status, &nullable.LastMessageSnippet, &nullable.LastMessageSenderID, &conversation.CreatedAt,
		&conversation.UpdatedAt, &conversation.ClosedAt); sErr != nil {
		return service.Conversation{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
	}

	if nullable.AssignedSupportID.Valid {
		conversation.AssignedSupportID = types.ID(nullable.AssignedSupportID.String)
	}

	if nullable.LastMessageSnippet.Valid {
		conversation.LastMessageSnippet = nullable.LastMessageSnippet.String
	}

	if nullable.LastMessageSenderID.Valid {
		conversation.LastMessageSenderID = types.ID(nullable.LastMessageSenderID.String)
	}

	return conversation, nil
}
