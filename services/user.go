package services

import (
	"context"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/token"
	"github.com/dato7898/grpc-tube/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken: accessToken,
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
	}, nil
}

func (s *Server) Me(ctx context.Context, req *emptypb.Empty) (*pb.AuthPayload, error) {
	payload := ctx.Value(AuthorizationPayloadKey).(*token.Payload)
	return &pb.AuthPayload{
		Id:        payload.ID.String(),
		Username:  payload.Username,
		IssuedAt:  timestamppb.New(payload.IssuedAt),
		ExpiredAt: timestamppb.New(payload.ExpiredAt),
	}, nil
}
