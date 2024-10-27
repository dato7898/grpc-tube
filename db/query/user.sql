-- name: CreateUser :one
INSERT INTO "user" (
    username,
    hashed_password,
    email
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: VerifyUser :exec
UPDATE "user"
SET verified = true
WHERE id = $1;