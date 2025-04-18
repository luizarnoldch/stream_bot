-- name: CreateUser :one
INSERT    INTO auth."user" (NAME, email)
VALUES    ($1, $2)
RETURNING id,
          NAME,
          email,
          created_at,
          updated_at;

-- name: GetUserByID :one
SELECT    id,
          NAME,
          email,
          created_at,
          updated_at
FROM      auth."user"
WHERE     id = $1;

-- name: ListUsers :many
SELECT    id,
          NAME,
          email,
          created_at,
          updated_at
FROM      auth."user"
ORDER BY  id;