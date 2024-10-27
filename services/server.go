package services

import (
	"log"
	"net"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server serves GRPC request for our service.
type Server struct {
	pb.UserServer
	store      db.Store
	grpcServer *grpc.Server
}

// NewServer creates a new GRPC server and register all servers
func NewServer(store db.Store) (*Server, error) {
	s := grpc.NewServer()

	server := &Server{
		store:      store,
		grpcServer: s,
	}

	pb.RegisterUserServer(s, server)
	reflection.Register(s)

	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start() error {
	port := "8080"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	return server.grpcServer.Serve(lis)
}
