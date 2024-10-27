package services

import (
	"context"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.Password,
		Email:          "email",
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
