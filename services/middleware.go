package services

import (
	"context"

	"github.com/dato7898/grpc-tube/pb"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const AuthorizationPayloadKey = contextKey("authorization_payload")

func (s *Server) authenticator(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	payload, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}
	ctx = context.WithValue(ctx, AuthorizationPayloadKey, payload)
	return ctx, nil
}

func authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	if pb.User_ServiceDesc.ServiceName == callMeta.Service {
		if pb.User_ServiceDesc.Methods[2].MethodName == callMeta.Method {
			return true
		}
	}
	return false
}
