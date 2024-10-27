package services

import (
	"context"
	"fmt"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/mail"
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

	code := util.RandomCode(6)

	tokenArg := db.CreateAuthTokenParams{
		UserID: user.ID,
		Code:   code,
	}

	s.store.CreateAuthToken(ctx, tokenArg)

	mail.Send([]string{req.Email}, fmt.Sprintf("Verification code: %s", code))

	return &pb.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *Server) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	authToken, err := s.store.GetLastAuthTokenByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if req.Code != authToken.Code {
		return nil, fmt.Errorf("invalid code")
	}

	err = s.store.VerifyUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyResponse{
		Success: true,
	}, nil
}
