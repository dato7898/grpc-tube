package services

import (
	"fmt"
	"log"
	"net"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/token"
	"github.com/dato7898/grpc-tube/util"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server serves GRPC request for our service.
type Server struct {
	pb.UserServer
	pb.VideoServer
	store      db.Store
	grpcServer *grpc.Server
	tokenMaker token.Maker
	config     util.Config
}

// NewServer creates a new GRPC server and register all servers
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
			selector.UnaryServerInterceptor(
				auth.UnaryServerInterceptor(server.authenticator),
				selector.MatchFunc(authMatcher),
			),
		),
		grpc.ChainStreamInterceptor(
			// Stream interceptors chain for streaming RPCs
			selector.StreamServerInterceptor(
				auth.StreamServerInterceptor(server.authenticator),
				selector.MatchFunc(authMatcher),
			),
		),
	)

	server.grpcServer = s

	pb.RegisterUserServer(s, server)
	pb.RegisterVideoServer(s, server)
	reflection.Register(s)

	return server, nil
}

// Start runs the GRPC server on a specific address
func (server *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	return server.grpcServer.Serve(lis)
}
