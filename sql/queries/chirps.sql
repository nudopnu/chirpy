-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetAllChirps :many
SELECT * from chirps 
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * from chirps
WHERE id=$1;

-- name: DeleteChirpById :exec
DELETE FROM chirps WHERE id=$1;

-- name: GetChirpsFromUser :many
SELECT * FROM chirps
WHERE user_id = $1;