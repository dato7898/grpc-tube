-- name: CreateVideo :one
INSERT INTO "video" (
    id,
    title,
    description
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAll :many
SELECT * FROM "video"
LIMIT $1
OFFSET $2;