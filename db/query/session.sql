-- name: CreateSession :one
INSERT INTO "sessions" (
    id,
    username,
    refresh_token,
    is_blocked,
    user_agent,
    client_ip,
    expires_at,
    created_at
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: GetSession :one
SELECT *
FROM "sessions"
WHERE id = $1
LIMIT 1;