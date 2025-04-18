// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: 01000_auth.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT    INTO auth."user" (NAME, phone)
VALUES    ($1, $2)
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at
`

func (q *Queries) CreateUser(ctx context.Context, name string, phone string) (AuthUser, error) {
	row := q.db.QueryRow(ctx, createUser, name, phone)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
DELETE    FROM auth."user"
WHERE     id = $1
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) (AuthUser, error) {
	row := q.db.QueryRow(ctx, deleteUser, id)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
WHERE     id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (AuthUser, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
WHERE     phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone string) (AuthUser, error) {
	row := q.db.QueryRow(ctx, getUserByPhone, phone)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT    id,
          NAME,
          phone,
          created_at,
          updated_at
FROM      auth."user"
ORDER BY  id
`

func (q *Queries) ListUsers(ctx context.Context) ([]AuthUser, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuthUser
	for rows.Next() {
		var i AuthUser
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Phone,
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

const updateUser = `-- name: UpdateUser :one
UPDATE    auth."user"
SET       NAME = COALESCE($2, NAME),
          phone = COALESCE($3, phone),
          updated_at = NOW()
WHERE     id = $1
RETURNING id,
          NAME,
          phone,
          created_at,
          updated_at
`

func (q *Queries) UpdateUser(ctx context.Context, iD int32, name string, phone string) (AuthUser, error) {
	row := q.db.QueryRow(ctx, updateUser, iD, name, phone)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
