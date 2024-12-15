package services

import (
	"fmt"
	"log"
	"net"
	"net/http"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/pb"
	"github.com/dato7898/grpc-tube/token"
	"github.com/dato7898/grpc-tube/util"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
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
	httpServer *http.Server
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

	grpcServer := grpc.NewServer(
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

	server.grpcServer = grpcServer

	pb.RegisterUserServer(grpcServer, server)
	pb.RegisterVideoServer(grpcServer, server)
	reflection.Register(grpcServer)

	grpcWebServer := grpcweb.WrapServer(grpcServer)

	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request method: %s", r.Method)

			// Добавляем CORS заголовки
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-Requested-With, grpc-status, grpc-message, X-User-Agent, X-Grpc-Web, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message, application/grpc-web-text, X-Grpc-Web, X-User-Agent, grpc-web-javascript/0.1")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if grpcWebServer.IsGrpcWebRequest(r) {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
		Addr: ":" + config.HttpServerPort,
	}

	server.httpServer = httpServer

	return server, nil
}

// Start runs the GRPC server on a specific address
func (server *Server) Start() error {

	go func() {
		lis, err := net.Listen("tcp", ":"+server.config.GrpcServerPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Pure gRPC server listening at %v", lis.Addr())
		if err := server.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	log.Printf("Starting gRPC-Web server on port :%v", server.config.HttpServerPort)
	err := server.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	return nil
}
