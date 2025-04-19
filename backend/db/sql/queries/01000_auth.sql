-- name: CreateUser :one
INSERT    INTO auth."user" (NAME, phone)
VALUES    ($1, $2)
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at;

-- name: GetUserByID :one
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
WHERE     id = $1;

-- name: GetUserByPhone :one
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
WHERE     phone = $1;

-- name: ListUsers :many
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
ORDER BY  id;

-- name: UpdateUser :one
UPDATE    auth."user"
SET       NAME = COALESCE($2, NAME),
          phone = COALESCE($3, phone),
          updated_at = NOW()
WHERE     id = $1
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at;

-- name: DeleteUser :one
DELETE    FROM auth."user"
WHERE     id = $1
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at;