package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryUpdateConversationSnippet = `
UPDATE conversations
SET last_message_snippet = $2,
    last_message_sender_id = $3
WHERE id = $1;`

func (d *DB) UpdateConversationSnippet(ctx context.Context, conversationID, lastMessageSenderID types.ID, snippet string) error {
	const op = "repository.postgres.update.UpdateConversationSnippet"

	cmdTag, err := d.conn.Conn().Exec(ctx, queryUpdateConversationSnippet,
		conversationID,
		snippet,
		lastMessageSenderID,
	)

	if err != nil {
		return richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected)
	}

	if cmdTag.RowsAffected() == 0 {
		// This might happen if the conversation ID is wrong, but it's not necessarily
		// a critical error that should stop the message flow. We can just log it.
		// For now, we return nil as the operation "succeeded" (no DB error).
		// Or return a specific "not found" error if the service layer needs to know.
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage("conversation not found for snippet update")
	}

	return nil
}

const queryAssignConversation = `UPDATE conversations
SET assigned_support_id = $1,
    status = 'open'
WHERE id = $2;`

func (d *DB) AssignConversation(ctx context.Context, conversationID, supportID types.ID) error {
	const op = "repository.postgres.update.AssignConversation"

	cmdTag, exErr := d.conn.Conn().Exec(ctx, queryAssignConversation, supportID, conversationID)
	if exErr != nil {
		return richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected)
	}

	if cmdTag.RowsAffected() == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage("conversation not found for assign")
	}

	return nil
}

const queryCloseConversation = `UPDATE conversations
SET status = 'closed'
WHERE id = $1;`

func (d *DB) CloseConversation(ctx context.Context, conversationID types.ID) error {
	const op = "repository.postgres.update.CloseConversation"

	cmdTag, exErr := d.conn.Conn().Exec(ctx, queryCloseConversation, conversationID)
	if exErr != nil {
		return richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected)
	}

	if cmdTag.RowsAffected() == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage("conversation not found for close")
	}

	return nil
}
