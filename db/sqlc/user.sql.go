// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
    username,
    hashed_password,
    email
) VALUES (
    $1, $2, $3
) RETURNING id, username, hashed_password, email, verified
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Email,
		&i.Verified,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, hashed_password, email, verified FROM "user"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Email,
		&i.Verified,
	)
	return i, err
}

const verifyUser = `-- name: VerifyUser :exec
UPDATE "user"
SET verified = true
WHERE id = $1
`

func (q *Queries) VerifyUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, verifyUser, id)
	return err
}
