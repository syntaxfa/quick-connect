package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

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

const querySaveMessage = `
INSERT INTO messages (id, conversation_id, sender_id, message_type, content, metadata, replied_to_message_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING created_at;` // Return created_at in case DB DEFAULT NOW() was used (though we pass it)

func (d *DB) SaveMessage(ctx context.Context, message service.Message) (service.Message, error) {
	const op = "repository.postgres.create.SaveMessage"

	var nullContent sql.NullString
	if message.Content != "" {
		nullContent.String = message.Content
		nullContent.Valid = true
	}

	var nullRepliedToID sql.NullString
	if message.RepliedToMessageID != "" {
		nullRepliedToID.String = string(message.RepliedToMessageID)
		nullRepliedToID.Valid = true
	}

	// Handle JSONB metadata
	var metaDataBytes []byte
	var err error

	if len(message.MetaData) > 0 {
		metaDataBytes, err = json.Marshal(message.MetaData)
		if err != nil {
			return service.Message{}, richerror.New(op).WithWrapError(err).WithKind(richerror.KindInvalid).
				WithMessage("failed to marshal metadata")
		}
	} else {
		metaDataBytes = nil // Use NULL for empty or nil metadata
	}

	if sErr := d.conn.Conn().QueryRow(ctx, querySaveMessage,
		message.ID,
		message.ConversationID,
		message.SenderID,
		message.MessageType,
		nullContent,
		metaDataBytes,
		nullRepliedToID,
		message.CreatedAt,
	).Scan(&message.CreatedAt); sErr != nil { // Scan the returned created_at back into the struct
		return service.Message{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected)
	}

	return message, nil
}
