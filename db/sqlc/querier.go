// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error)
	GetAll(ctx context.Context, arg GetAllParams) ([]Video, error)
	GetUser(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)
