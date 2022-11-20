-- name: CreateUser :one
INSERT INTO "users" (owner, username, salt, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM "users"
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "users"
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUserOwner :one
UPDATE "users"
SET owner = $2
WHERE id = $1
RETURNING id;

-- name: UpdateCredentials :exec
UPDATE "users"
SET username = $2,
    password = $3
WHERE id = $1;

-- name: DeleteUser :one
DELETE
FROM "users"
WHERE id = $1
RETURNING id;