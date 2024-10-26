package main

import (
	"log"
	"net"

	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/services/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := "1531"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServer(s, &user.Server{})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
