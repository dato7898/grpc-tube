-- name: CreateAuthToken :exec
INSERT INTO "auth_token" (
    user_id,
    code
) VALUES (
    $1, $2
);

-- name: GetLastAuthTokenByUserId :one
SELECT * FROM "auth_token"
WHERE user_id = $1
ORDER BY expired_at DESC
LIMIT 1;
