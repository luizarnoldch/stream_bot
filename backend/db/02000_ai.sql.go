// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: 02000_ai.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createConversation = `-- name: CreateConversation :one
INSERT    INTO ai.conversation (NAME, user_id)
VALUES    ($1, $2)
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
`

func (q *Queries) CreateConversation(ctx context.Context, name string, userID int32) (AiConversation, error) {
	row := q.db.QueryRow(ctx, createConversation, name, userID)
	var i AiConversation
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createMessage = `-- name: CreateMessage :one
INSERT    INTO ai.message (CONTENT, ROLE, conversation_id)
VALUES    ($1, $2, $3)
RETURNING id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
`

func (q *Queries) CreateMessage(ctx context.Context, content string, role string, conversationID int32) (AiMessage, error) {
	row := q.db.QueryRow(ctx, createMessage, content, role, conversationID)
	var i AiMessage
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Role,
		&i.ConversationID,
		&i.Tokens,
		&i.FinishReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteConversation = `-- name: DeleteConversation :one
DELETE    FROM ai.conversation
WHERE     id = $1
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
`

func (q *Queries) DeleteConversation(ctx context.Context, id int32) (AiConversation, error) {
	row := q.db.QueryRow(ctx, deleteConversation, id)
	var i AiConversation
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMessage = `-- name: DeleteMessage :one
DELETE    FROM ai.message
WHERE     id = $1
RETURNING id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
`

func (q *Queries) DeleteMessage(ctx context.Context, id int32) (AiMessage, error) {
	row := q.db.QueryRow(ctx, deleteMessage, id)
	var i AiMessage
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Role,
		&i.ConversationID,
		&i.Tokens,
		&i.FinishReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getConversationByID = `-- name: GetConversationByID :one
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     id = $1
`

func (q *Queries) GetConversationByID(ctx context.Context, id int32) (AiConversation, error) {
	row := q.db.QueryRow(ctx, getConversationByID, id)
	var i AiConversation
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getConversationByUUID = `-- name: GetConversationByUUID :one
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     UUID = $1
`

func (q *Queries) GetConversationByUUID(ctx context.Context, uuid pgtype.UUID) (AiConversation, error) {
	row := q.db.QueryRow(ctx, getConversationByUUID, uuid)
	var i AiConversation
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMessageByID = `-- name: GetMessageByID :one
SELECT    id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
FROM      ai.message
WHERE     id = $1
`

func (q *Queries) GetMessageByID(ctx context.Context, id int32) (AiMessage, error) {
	row := q.db.QueryRow(ctx, getMessageByID, id)
	var i AiMessage
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Role,
		&i.ConversationID,
		&i.Tokens,
		&i.FinishReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listConversationsByUser = `-- name: ListConversationsByUser :many
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     user_id = $1
ORDER BY  created_at
`

func (q *Queries) ListConversationsByUser(ctx context.Context, userID int32) ([]AiConversation, error) {
	rows, err := q.db.Query(ctx, listConversationsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AiConversation
	for rows.Next() {
		var i AiConversation
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMessagesByConversation = `-- name: ListMessagesByConversation :many
SELECT    id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
FROM      ai.message
WHERE     conversation_id = $1
ORDER BY  created_at
`

func (q *Queries) ListMessagesByConversation(ctx context.Context, conversationID int32) ([]AiMessage, error) {
	rows, err := q.db.Query(ctx, listMessagesByConversation, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AiMessage
	for rows.Next() {
		var i AiMessage
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.Role,
			&i.ConversationID,
			&i.Tokens,
			&i.FinishReason,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateConversation = `-- name: UpdateConversation :one
UPDATE    ai.conversation
SET       NAME = COALESCE($2, NAME),
          updated_at = NOW()
WHERE     id = $1
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
`

func (q *Queries) UpdateConversation(ctx context.Context, iD int32, name string) (AiConversation, error) {
	row := q.db.QueryRow(ctx, updateConversation, iD, name)
	var i AiConversation
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateMessage = `-- name: UpdateMessage :one
UPDATE    ai.message
SET       CONTENT = COALESCE($2, CONTENT),
          tokens = COALESCE($3, tokens),
          finish_reason = COALESCE($4, finish_reason),
          updated_at = NOW()
WHERE     id = $1
RETURNING id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
`

type UpdateMessageParams struct {
	ID           int32       `json:"id"`
	Content      string      `json:"content"`
	Tokens       pgtype.Int4 `json:"tokens"`
	FinishReason pgtype.Text `json:"finish_reason"`
}

func (q *Queries) UpdateMessage(ctx context.Context, arg UpdateMessageParams) (AiMessage, error) {
	row := q.db.QueryRow(ctx, updateMessage,
		arg.ID,
		arg.Content,
		arg.Tokens,
		arg.FinishReason,
	)
	var i AiMessage
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Role,
		&i.ConversationID,
		&i.Tokens,
		&i.FinishReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
