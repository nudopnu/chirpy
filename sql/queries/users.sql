-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * from users
WHERE email=$1;

-- name: GetUserById :one
SELECT * from users
WHERE id=$1;

-- name: UpdateUser :one
UPDATE users
SET email=$2, hashed_password=$3, updated_at=NOW()
WHERE id=$1 RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red=TRUE
WHERE id=$1;