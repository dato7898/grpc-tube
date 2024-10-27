package services

import (
	"context"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/util"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
