-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
  AND is_deleted = false
LIMIT 1;
-- name: GetAccounts :many
SELECT *
FROM accounts
WHERE owner = $1
  AND is_deleted = false;
-- name: GetDeletedAccounts :many
SELECT *
FROM accounts
WHERE owner = $1
  AND is_deleted = true;
-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteAccount :exec
UPDATE accounts
SET is_deleted = true
WHERE id = $1;
-- name: RestoreAccount :exec
UPDATE accounts
SET is_deleted = false
WHERE id = $1;