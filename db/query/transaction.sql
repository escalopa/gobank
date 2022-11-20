-- name: CreateTransfer :one
INSERT INTO "transfers" (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountTransfer :one
SELECT *
FROM "transfers"
WHERE from_account_id = $1
LIMIT $2 OFFSET $3;

-- name: GetTransfer :one
SELECT *
FROM "transfers"
WHERE id = $1
LIMIT 1;

-- name: GetUserTransfers :many
SELECT a.id, a.currency, from_account_id, to_account_id, amount, reverted, t.created_at
FROM "transfers" t
         JOIN accounts a on a.id = t.from_account_id  AND a.id = t.to_account_id
WHERE a.user_id = $1
LIMIT $2 OFFSET $3;

-- name: GetUserSentTransfers :many
SELECT a.id, a.currency, from_account_id, to_account_id, amount, t.created_at
FROM "transfers" t
         JOIN accounts a on a.id = t.from_account_id  AND t.reverted = false
WHERE a.user_id = $1
LIMIT $2 OFFSET $3;

-- name: GetUserReceivedTransfers :many
SELECT a.id, a.currency, from_account_id, to_account_id, amount, t.created_at
FROM "transfers" t
         JOIN accounts a on a.id = t.to_account_id AND t.reverted = false
WHERE a.user_id = $1 ORDER BY t.id
LIMIT $2 OFFSET $3;

-- name: ListTransfers :many
SELECT *
FROM "transfers"
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteTransfer :one
UPDATE "transfers"
SET reverted = true
WHERE id = $1 RETURNING id;