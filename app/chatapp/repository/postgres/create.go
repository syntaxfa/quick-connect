package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryCreateActiveConversation = `INSERT INTO conversations (id, client_user_id, status)
VALUES ($1, $2, $3)
RETURNING id, client_user_id, status, created_at, updated_at;`

func (d *DB) CreateActiveConversation(ctx context.Context, id, userID types.ID,
	conversationStatus service.ConversationStatus) (service.Conversation, error) {
	const op = "repository.postgres.create.CreateActiveConversation"

	var conversation service.Conversation

	if sErr := d.conn.Conn().QueryRow(ctx, queryCreateActiveConversation, id, userID, conversationStatus).Scan(
		&conversation.ID, &conversation.ClientUserID, &conversation.Status, &conversation.CreatedAt, &conversation.UpdatedAt); sErr != nil {
		return service.Conversation{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
	}

	return conversation, nil
}

func (d *DB) SaveMessage(_ service.Message) error {
	// TODO implement me
	panic("implement me")
}
