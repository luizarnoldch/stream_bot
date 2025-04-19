-- name: CreateConversation :one
INSERT    INTO ai.conversation (NAME, user_id)
VALUES    ($1, $2)
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at;

-- name: GetConversationByID :one
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     id = $1;

-- name: GetConversationByUUID :one
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     UUID = $1;

-- name: ListConversationsByUser :many
SELECT    id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at
FROM      ai.conversation
WHERE     user_id = $1
ORDER BY  created_at;

-- name: UpdateConversation :one
UPDATE    ai.conversation
SET       NAME = COALESCE($2, NAME),
          updated_at = NOW()
WHERE     id = $1
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at;

-- name: DeleteConversation :one
DELETE    FROM ai.conversation
WHERE     id = $1
RETURNING id,
          UUID,
          NAME,
          user_id,
          created_at,
          updated_at;

-- name: CreateMessage :one
INSERT    INTO ai.message (CONTENT, ROLE, conversation_id)
VALUES    ($1, $2, $3)
RETURNING id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at;

-- name: GetMessageByID :one
SELECT    id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at
FROM      ai.message
WHERE     id = $1;

-- name: ListMessagesByConversation :many
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
ORDER BY  created_at;

-- name: UpdateMessage :one
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
          updated_at;

-- name: DeleteMessage :one
DELETE    FROM ai.message
WHERE     id = $1
RETURNING id,
          CONTENT,
          ROLE,
          conversation_id,
          tokens,
          finish_reason,
          created_at,
          updated_at;