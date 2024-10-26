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

-- name: GetChirpsFiltered :many
SELECT * FROM chirps
WHERE ($1::UUID = '00000000-0000-0000-0000-000000000000'::UUID OR user_id = $1::UUID)
ORDER BY
    CASE WHEN $2::TEXT = 'asc' THEN created_at END ASC,
    CASE WHEN $2::TEXT = 'desc' THEN created_at END DESC;