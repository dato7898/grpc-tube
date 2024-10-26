package user

import (
	"context"

	"github.com/dato7898/grpc-tube/pb"
)

type Server struct {
	pb.UserServer
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{
		Id:       1,
		Username: req.Username,
	}, nil
}
