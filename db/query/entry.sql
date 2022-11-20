-- name: CreateEntry :one
INSERT INTO "entries" (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccountEntry :one
SELECT *
FROM "entries"
WHERE account_id = $1
LIMIT 1;

-- name: GetEntry :one
SELECT *
FROM "entries"
WHERE id = $1
LIMIT 1;

-- name: GetUserEntries :many
SELECT account_id, amount, currency, e.created_at
FROM "entries" e
         JOIN accounts a on a.id = e.account_id
WHERE a.user_id = $1
LIMIT $2 OFFSET $3;

-- name: ListEntries :many
SELECT *
FROM "entries"
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteEntry :exec
DELETE
FROM "entries"
WHERE id = $1;