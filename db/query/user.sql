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
WHERE username = $1 LIMIT 1;