-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, created_at, updated_at) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * from refresh_tokens WHERE token = $1;

-- name: GetUserByRefreshToken :one
SELECT users.* from refresh_tokens
JOIN users ON users.id = refresh_tokens.user_id
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at=$2, updated_at=$2
WHERE token=$1;