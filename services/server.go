package services

import (
	"context"
	"fmt"
	"log"
	"net"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/token"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type contextKey string

const AuthorizationPayloadKey = contextKey("authorization_payload")

// Server serves GRPC request for our service.
type Server struct {
	pb.UserServer
	store      db.Store
	grpcServer *grpc.Server
	tokenMaker token.Maker
}

// NewServer creates a new GRPC server and register all servers
func NewServer(store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
		selector.UnaryServerInterceptor(
			auth.UnaryServerInterceptor(server.authenticator),
			selector.MatchFunc(authMatcher),
		),
	))

	server.grpcServer = s

	pb.RegisterUserServer(s, server)
	reflection.Register(s)

	return server, nil
}

// Start runs the GRPC server on a specific address
func (server *Server) Start() error {
	port := "8080"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	return server.grpcServer.Serve(lis)
}

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
