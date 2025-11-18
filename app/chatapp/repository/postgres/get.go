package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const queryGetUserActiveConversation = `SELECT id, client_user_id, assigned_support_id, status, last_message_snippet,
last_message_sender_id, created_at, updated_at, closed_at
FROM conversations
WHERE client_user_id = $1 AND status IN ('new', 'open', 'bot_handling')
limit 1;`

const queryGetConversationByID = `SELECT id, client_user_id, assigned_support_id, status, last_message_snippet,
last_message_sender_id, created_at, updated_at, closed_at
FROM conversations
WHERE id = $1
limit 1;`

type nullableFields struct {
	AssignedSupportID   sql.NullString
	LastMessageSnippet  sql.NullString
	LastMessageSenderID sql.NullString
	ClosedAt            pq.NullTime
}

func (d *DB) GetUserActiveConversation(ctx context.Context, userID types.ID) (service.Conversation, error) {
	const op = "repository.postgres.get.GetUserActiveConversation"

	return d.GetConversationBy(ctx, op, queryGetUserActiveConversation, userID)
}

func (d *DB) GetConversationByID(ctx context.Context, conversationID types.ID) (service.Conversation, error) {
	const op = "repository.postgres.get.GetConversationByID"

	return d.GetConversationBy(ctx, op, queryGetConversationByID, conversationID)
}

func (d *DB) GetConversationBy(ctx context.Context, op string, query string, arg interface{}) (service.Conversation, error) {
	var conversation service.Conversation
	var nullable nullableFields

	if sErr := d.conn.Conn().QueryRow(ctx, query, arg).Scan(&conversation.ID, &conversation.ClientUserID, &nullable.AssignedSupportID,
		&conversation.Status, &nullable.LastMessageSnippet, &nullable.LastMessageSenderID, &conversation.CreatedAt,
		&conversation.UpdatedAt, &nullable.ClosedAt); sErr != nil {
		if errors.Is(sErr, pgx.ErrNoRows) {
			return service.Conversation{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindNotFound)
		}

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

	if nullable.ClosedAt.Valid {
		conversation.ClosedAt = &nullable.ClosedAt.Time
	}

	return conversation, nil
}

// buildConversationListFilters is a helper to build dynamic WHERE clauses.
func (d *DB) buildConversationListFilters(assignedSupportID types.ID, statuses []service.ConversationStatus) (string, []interface{}) {
	args := make([]interface{}, 0)
	whereClauses := make([]string, 0)
	argCount := 1

	if assignedSupportID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("assigned_support_id = $%d", argCount))
		args = append(args, assignedSupportID)
		argCount++
	}

	if len(statuses) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("status = ANY($%d)", argCount))
		// pq.Array is used to pass a Go slice as a PostgreSQL array
		args = append(args, pq.Array(statuses))
		// argCount++ // pq.Array counts as one argument
	}

	whereQuery := ""
	if len(whereClauses) > 0 {
		whereQuery = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	return whereQuery, args
}

// GetConversationList retrieves a paginated list of conversations.
func (d *DB) GetConversationList(ctx context.Context, paginated paginate.RequestBase,
	assignedSupportID types.ID, statuses []service.ConversationStatus) ([]service.Conversation, paginate.ResponseBase, error) {
	const op = "repository.postgres.get.GetConversationList"

	offset := (paginated.CurrentPage - 1) * paginated.PageSize
	limit := paginated.PageSize

	// Default sorting for conversations is usually by update time
	sortColumn := "id"
	sortDirection := "ASC"
	if paginated.Descending {
		sortDirection = "DESC"
	}

	whereQuery, args := d.buildConversationListFilters(assignedSupportID, statuses)

	const (
		limitArgDelta  = 1
		offsetArgDelta = 2
	)

	argCount := len(args)
	query := fmt.Sprintf(`
       SELECT id, client_user_id, assigned_support_id, status,
              last_message_snippet, last_message_sender_id,
              created_at, updated_at, closed_at
       FROM conversations
       %s
       ORDER BY %s %s
       LIMIT $%d OFFSET $%d`,
		whereQuery, sortColumn, sortDirection, argCount+limitArgDelta, argCount+offsetArgDelta)

	mainQueryArgs := make([]interface{}, 0, len(args)+offsetArgDelta)
	mainQueryArgs = append(mainQueryArgs, args...)
	mainQueryArgs = append(mainQueryArgs, limit, offset)

	rows, qErr := d.conn.Conn().Query(ctx, query, mainQueryArgs...)
	if qErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(qErr).WithKind(richerror.KindUnexpected).
			WithMessage("conversation list query error")
	}
	defer rows.Close()

	var conversations []service.Conversation
	for rows.Next() {
		var conv service.Conversation
		var nullable nullableFields

		if sErr := rows.Scan(
			&conv.ID, &conv.ClientUserID, &nullable.AssignedSupportID, &conv.Status,
			&nullable.LastMessageSnippet, &nullable.LastMessageSenderID,
			&conv.CreatedAt, &conv.UpdatedAt, &nullable.ClosedAt,
		); sErr != nil {
			return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
				WithMessage("scan error")
		}

		// Handle nullable fields
		if nullable.AssignedSupportID.Valid {
			conv.AssignedSupportID = types.ID(nullable.AssignedSupportID.String)
		}
		if nullable.LastMessageSnippet.Valid {
			conv.LastMessageSnippet = nullable.LastMessageSnippet.String
		}
		if nullable.LastMessageSenderID.Valid {
			conv.LastMessageSenderID = types.ID(nullable.LastMessageSenderID.String)
		}
		if nullable.ClosedAt.Valid {
			conv.ClosedAt = &nullable.ClosedAt.Time
		}

		conversations = append(conversations, conv)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected)
	}

	// --- Get Total Count ---
	countQuery := fmt.Sprintf(`
       SELECT COUNT(*)
       FROM conversations
       %s`,
		whereQuery)

	var totalCount uint64
	if sErr := d.conn.Conn().QueryRow(ctx, countQuery, args...).Scan(&totalCount); sErr != nil {
		return nil, paginate.ResponseBase{}, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected).
			WithMessage("count query error")
	}

	return conversations, paginate.ResponseBase{
		CurrentPage:  paginated.CurrentPage,
		PageSize:     paginated.PageSize,
		TotalNumbers: totalCount,
		TotalPage:    (totalCount + paginated.PageSize - 1) / paginated.PageSize,
	}, nil
}

const queryGetConversationParticipants = `
SELECT client_user_id, assigned_support_id
FROM conversations
WHERE id = $1;`

func (d *DB) GetConversationParticipants(ctx context.Context, conversationID types.ID) ([]types.ID, error) {
	const op = "repository.postgres.get.GetConversationParticipants"

	var clientUserID types.ID
	var nullAssignedSupportID sql.NullString

	if err := d.conn.Conn().QueryRow(ctx, queryGetConversationParticipants, conversationID).Scan(
		&clientUserID,
		&nullAssignedSupportID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, richerror.New(op).WithWrapError(err).WithKind(richerror.KindNotFound)
		}

		return nil, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected)
	}

	participants := []types.ID{clientUserID}

	if nullAssignedSupportID.Valid && nullAssignedSupportID.String != "" {
		participants = append(participants, types.ID(nullAssignedSupportID.String))
	}

	return participants, nil
}
