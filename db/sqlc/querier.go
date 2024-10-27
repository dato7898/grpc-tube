// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateAuthToken(ctx context.Context, arg CreateAuthTokenParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetLastAuthTokenByUserId(ctx context.Context, userID int64) (AuthToken, error)
	GetUser(ctx context.Context, id int64) (User, error)
	VerifyUser(ctx context.Context, id int64) error
}

var _ Querier = (*Queries)(nil)
