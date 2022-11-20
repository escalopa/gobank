// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO "accounts" (user_id, balance, currency)
VALUES ($1, 0, $2)
RETURNING id, user_id, balance, currency, created_at
`

type CreateAccountParams struct {
	UserID   int64    `json:"user_id"`
	Currency Currency `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.UserID, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :one
DELETE
FROM "accounts"
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteAccount, id)
	err := row.Scan(&id)
	return id, err
}

const getAccount = `-- name: GetAccount :one
SELECT id, user_id, balance, currency, created_at
FROM "accounts"
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getUserAccount = `-- name: GetUserAccount :one
SELECT id, user_id, balance, currency, created_at
FROM "accounts"
WHERE user_id = $1
  AND currency = $2
LIMIT 1
`

type GetUserAccountParams struct {
	UserID   int64    `json:"user_id"`
	Currency Currency `json:"currency"`
}

func (q *Queries) GetUserAccount(ctx context.Context, arg GetUserAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, getUserAccount, arg.UserID, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getUserAccounts = `-- name: GetUserAccounts :many
SELECT id, user_id, balance, currency, created_at
FROM "accounts"
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3
`

type GetUserAccountsParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUserAccounts(ctx context.Context, arg GetUserAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, getUserAccounts, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Balance,
			&i.Currency,
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

const listAccounts = `-- name: ListAccounts :many
SELECT id, user_id, balance, currency, created_at
FROM "accounts"
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Balance,
			&i.Currency,
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

const updateBalance = `-- name: UpdateBalance :one
UPDATE "accounts"
SET balance = balance + $2
WHERE id = $1
RETURNING balance
`

type UpdateBalanceParams struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (float64, error) {
	row := q.db.QueryRowContext(ctx, updateBalance, arg.ID, arg.Balance)
	var balance float64
	err := row.Scan(&balance)
	return balance, err
}
