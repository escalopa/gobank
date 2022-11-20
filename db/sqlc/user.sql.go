// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" (owner, username, salt, password)
VALUES ($1, $2, $3, $4)
RETURNING id, owner, username, salt, password, created_at
`

type CreateUserParams struct {
	Owner    string `json:"owner"`
	Username string `json:"username"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Owner,
		arg.Username,
		arg.Salt,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Username,
		&i.Salt,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
DELETE
FROM "users"
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteUser, id)
	err := row.Scan(&id)
	return id, err
}

const getUser = `-- name: GetUser :one
SELECT id, owner, username, salt, password, created_at
FROM "users"
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Username,
		&i.Salt,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, owner, username, salt, password, created_at
FROM "users"
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Username,
			&i.Salt,
			&i.Password,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCredentials = `-- name: UpdateCredentials :exec
UPDATE "users"
SET username = $2,
    password = $3
WHERE id = $1
`

type UpdateCredentialsParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) UpdateCredentials(ctx context.Context, arg UpdateCredentialsParams) error {
	_, err := q.db.ExecContext(ctx, updateCredentials, arg.ID, arg.Username, arg.Password)
	return err
}

const updateUserOwner = `-- name: UpdateUserOwner :one
UPDATE "users"
SET owner = $2
WHERE id = $1
RETURNING id
`

type UpdateUserOwnerParams struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
}

func (q *Queries) UpdateUserOwner(ctx context.Context, arg UpdateUserOwnerParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateUserOwner, arg.ID, arg.Owner)
	var id int64
	err := row.Scan(&id)
	return id, err
}
