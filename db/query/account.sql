-- name: CreateAccount :one
INSERT INTO "accounts" (user_id, balance, currency)
VALUES ($1, 0, $2)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM "accounts"
WHERE id = $1
LIMIT 1;

-- name: GetUserAccount :one
SELECT *
FROM "accounts"
WHERE user_id = $1
  AND currency = $2
LIMIT 1;

-- name: GetUserAccounts :many
SELECT *
FROM "accounts"
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: ListAccounts :many
SELECT *
FROM "accounts"
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateBalance :one
UPDATE "accounts"
SET balance = balance + $2
WHERE id = $1
RETURNING balance;

-- name: DeleteAccount :one
DELETE
FROM "accounts"
WHERE id = $1
RETURNING id;